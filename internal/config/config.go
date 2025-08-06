package config

import (
	"fmt"
	"strings"
	"time"
)

// Config holds the server configuration
type Config struct {
	// AllowList contains allowed orgs/users and repos
	// Entries can be:
	//   - "org" or "user" to allow all repos from that org/user
	//   - "owner/repo" to allow a specific repo
	AllowList []string

	// HTTPCacheDuration is how long browsers/CDNs should cache banner images
	HTTPCacheDuration time.Duration

	// APICacheDuration is how long to cache GitHub API responses
	APICacheDuration time.Duration
}

// NewConfig creates a new config from a list of allowed entries
func NewConfig(allowList []string) *Config {
	return &Config{
		AllowList:         allowList,
		HTTPCacheDuration: 1 * time.Hour, // Default to 1 hour
		APICacheDuration:  1 * time.Hour, // Default to 1 hour
	}
}

// IsAllowed checks if a repository is allowed
func (c *Config) IsAllowed(owner, repo string) bool {
	// Empty allowlist means allow everything
	if len(c.AllowList) == 0 {
		return true
	}

	fullRepo := fmt.Sprintf("%s/%s", owner, repo)

	// Check each entry in the allowlist
	for _, allowed := range c.AllowList {
		// Check if it's a specific repo match
		if strings.EqualFold(allowed, fullRepo) {
			return true
		}

		// Check if it's an org/user match (no slash means org/user)
		if !strings.Contains(allowed, "/") && strings.EqualFold(allowed, owner) {
			return true
		}
	}

	return false
}
