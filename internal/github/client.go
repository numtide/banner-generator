package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
	"time"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	cache  map[string]*cacheEntry
}

type cacheEntry struct {
	data      *Repository
	timestamp time.Time
}

// NewClient creates a new GitHub client
func NewClient(token string) *Client {
	ctx := context.Background()
	var tc *oauth2.TokenSource

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = &ts
	}

	var ghClient *github.Client
	if tc != nil {
		httpClient := oauth2.NewClient(ctx, *tc)
		ghClient = github.NewClient(httpClient)
	} else {
		ghClient = github.NewClient(nil)
	}

	return &Client{
		client: ghClient,
		cache:  make(map[string]*cacheEntry),
	}
}

// GetRepositoryData fetches repository metadata from GitHub
func (c *Client) GetRepositoryData(ctx context.Context, owner, repo string) (*Repository, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("%s/%s", owner, repo)
	if entry, ok := c.cache[cacheKey]; ok {
		// Cache entries are valid for 5 minutes
		if time.Since(entry.timestamp) < 5*time.Minute {
			return entry.data, nil
		}
	}

	// Fetch repository information
	repository, _, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}

	// Use the original repo name as provided by the user to preserve capitalization
	data := &Repository{
		Name:            repo,
		Description:     repository.GetDescription(),
		Owner:           owner,
		Language:        repository.GetLanguage(),
		StargazersCount: repository.GetStargazersCount(),
		ForksCount:      repository.GetForksCount(),
	}

	// Update cache
	c.cache[cacheKey] = &cacheEntry{
		data:      data,
		timestamp: time.Now(),
	}

	return data, nil
}

// ValidateRepository checks if a repository exists and is accessible
func (c *Client) ValidateRepository(ctx context.Context, owner, repo string) error {
	_, _, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("repository validation failed: %w", err)
	}
	return nil
}
