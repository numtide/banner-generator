package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// AppConfig represents the complete application configuration
type AppConfig struct {
	// Server configuration
	Server ServerConfig `toml:"server"`

	// Font configuration
	Fonts FontsConfig `toml:"fonts"`

	// Template configuration
	Templates TemplateConfig `toml:"templates"`

	// GitHub configuration
	GitHub GitHubConfig `toml:"github"`

	// Access control configuration
	AccessControl AccessControlConfig `toml:"access_control"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Port         int    `toml:"port"`
	Host         string `toml:"host"`
	ReadTimeout  string `toml:"read_timeout"`
	WriteTimeout string `toml:"write_timeout"`
}

// FontsConfig contains font-related settings
type FontsConfig struct {
	// Path to fonts directory
	FontsDir string `toml:"fonts_dir"`

	// Path to web fonts directory (for serving web fonts)
	WebFontsDir string `toml:"web_fonts_dir"`

	// Default font family
	DefaultFamily string `toml:"default_family"`

	// Enable web fonts in SVG output
	EnableWebFonts bool `toml:"enable_web_fonts"`

	// Base URL for web fonts (when serving fonts via HTTP)
	WebFontsBaseURL string `toml:"web_fonts_base_url"`
}

// TemplateConfig contains template settings
type TemplateConfig struct {
	// Path to templates directory
	TemplatesDir string `toml:"templates_dir"`

	// Default template to use
	DefaultTemplate string `toml:"default_template"`

	// Available templates
	Templates map[string]TemplateEntry `toml:"templates"`
}

// TemplateEntry represents a single template
type TemplateEntry struct {
	Path        string `toml:"path"`
	Description string `toml:"description"`
}

// GitHubConfig contains GitHub API settings
type GitHubConfig struct {
	// GitHub API token (can be overridden by env var)
	Token string `toml:"token,omitempty"`
}

// AccessControlConfig contains access control settings
type AccessControlConfig struct {
	// Enable access control
	Enabled bool `toml:"enabled"`

	// Allowed GitHub organizations
	AllowedOrgs []string `toml:"allowed_orgs"`

	// Allowed GitHub users
	AllowedUsers []string `toml:"allowed_users"`
}

// LoadConfig loads configuration from a TOML file
func LoadConfig(path string) (*AppConfig, error) {
	// Start with default configuration
	config := DefaultConfig()

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil // Return defaults if file doesn't exist
	}

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse TOML
	if err := toml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	config.applyEnvOverrides()

	// Resolve relative paths
	if err := config.resolvePaths(filepath.Dir(path)); err != nil {
		return nil, fmt.Errorf("failed to resolve paths: %w", err)
	}

	return config, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Port:         8080,
			Host:         "0.0.0.0",
			ReadTimeout:  "30s",
			WriteTimeout: "30s",
		},
		Fonts: FontsConfig{
			FontsDir:        "deploy/fonts",
			WebFontsDir:     "deploy/fonts/web",
			DefaultFamily:   "GT Pressura",
			EnableWebFonts:  false,
			WebFontsBaseURL: "",
		},
		Templates: TemplateConfig{
			TemplatesDir:    "deploy/templates",
			DefaultTemplate: "banner",
			Templates: map[string]TemplateEntry{
				"banner": {
					Path:        "banner.svg.mustache",
					Description: "Default banner template",
				},
				"banner-stats": {
					Path:        "banner-stats.svg.mustache",
					Description: "Banner with repository statistics",
				},
			},
		},
		GitHub: GitHubConfig{
			Token: "",
		},
		AccessControl: AccessControlConfig{
			Enabled:      false,
			AllowedOrgs:  []string{},
			AllowedUsers: []string{},
		},
	}
}

// applyEnvOverrides applies environment variable overrides
func (c *AppConfig) applyEnvOverrides() {
	// Server
	if port := os.Getenv("PORT"); port != "" {
		if _, err := fmt.Sscanf(port, "%d", &c.Server.Port); err != nil {
			// Log error but continue with default port
			fmt.Fprintf(os.Stderr, "warning: invalid PORT value '%s', using default: %v\n", port, err)
		}
	}

	// GitHub
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		c.GitHub.Token = token
	}

	// Access Control
	if enabled := os.Getenv("ACCESS_CONTROL_ENABLED"); enabled == "true" {
		c.AccessControl.Enabled = true
	}
}

// resolvePaths resolves relative paths in the configuration
func (c *AppConfig) resolvePaths(basePath string) error {
	// Resolve fonts directory
	if !filepath.IsAbs(c.Fonts.FontsDir) {
		c.Fonts.FontsDir = filepath.Join(basePath, c.Fonts.FontsDir)
	}
	if !filepath.IsAbs(c.Fonts.WebFontsDir) {
		c.Fonts.WebFontsDir = filepath.Join(basePath, c.Fonts.WebFontsDir)
	}

	// Resolve template directory
	if !filepath.IsAbs(c.Templates.TemplatesDir) {
		c.Templates.TemplatesDir = filepath.Join(basePath, c.Templates.TemplatesDir)
	}

	return nil
}
