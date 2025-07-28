package cli

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

// Uploader handles PNG upload to GitHub repositories
type Uploader struct {
	client *github.Client
}

// NewUploader creates a new GitHub uploader
func NewUploader(token string) *Uploader {
	if token == "" {
		fmt.Println("Warning: No GitHub token provided. Upload functionality will be limited.")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Uploader{
		client: github.NewClient(tc),
	}
}

// UploadBanner uploads a PNG file to a GitHub repository
func (u *Uploader) UploadBanner(repoPath, pngPath string) error {
	return u.UploadBannerWithOptions(repoPath, pngPath, "main", "assets/banner.png", "Update repository banner")
}

// UploadBannerWithOptions uploads a PNG file with custom options
func (u *Uploader) UploadBannerWithOptions(repoPath, pngPath, branch, targetPath, commitMsg string) error {
	// Parse owner/repo format
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format, expected owner/repo")
	}
	owner, repo := parts[0], parts[1]

	// Read PNG file
	fmt.Printf("Reading PNG file: %s\n", pngPath)
	content, err := os.ReadFile(pngPath)
	if err != nil {
		return fmt.Errorf("failed to read PNG file: %w", err)
	}

	// Encode content as base64
	encodedContent := base64.StdEncoding.EncodeToString(content)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Check if file already exists
	fmt.Printf("Checking if %s exists in %s/%s...\n", targetPath, owner, repo)
	fileContent, _, _, err := u.client.Repositories.GetContents(ctx, owner, repo, targetPath, &github.RepositoryContentGetOptions{
		Ref: branch,
	})

	var fileSHA *string
	if err == nil && fileContent != nil {
		fileSHA = fileContent.SHA
		fmt.Println("File exists, will update it.")
	} else {
		fmt.Println("File doesn't exist, will create it.")
	}

	// Create or update file
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(commitMsg),
		Content: []byte(encodedContent),
		Branch:  github.String(branch),
	}

	if fileSHA != nil {
		opts.SHA = fileSHA
	}

	fmt.Printf("Uploading banner to %s/%s at %s...\n", owner, repo, targetPath)
	_, _, err = u.client.Repositories.CreateFile(ctx, owner, repo, targetPath, opts)
	if err != nil {
		// If create fails and we have a SHA, try update
		if fileSHA != nil {
			_, _, err = u.client.Repositories.UpdateFile(ctx, owner, repo, targetPath, opts)
		}
		if err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}
	}

	fmt.Printf("Successfully uploaded banner to: https://github.com/%s/%s/blob/%s/%s\n",
		owner, repo, branch, targetPath)
	return nil
}
