package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// ConfigLoader handles loading and merging configuration from multiple sources
type ConfigLoader struct {
	configPaths []string
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		configPaths: getDefaultConfigPaths(),
	}
}

// getDefaultConfigPaths returns the default paths to search for config files
func getDefaultConfigPaths() []string {
	paths := []string{
		"deploy/banner-generator.toml",
		"banner-generator.toml",
		"config.toml",
		".banner-generator.toml",
	}

	// Add user config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		paths = append(paths, filepath.Join(configDir, "banner-generator", "config.toml"))
	}

	// Add system config directory
	paths = append(paths, "/etc/banner-generator/config.toml")

	return paths
}

// LoadConfig loads the application configuration
func (l *ConfigLoader) LoadConfig(explicitPath string) (*AppConfig, error) {
	var config *AppConfig
	var configPath string

	// If explicit path is provided, use only that
	if explicitPath != "" {
		config, err := LoadConfig(explicitPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", explicitPath, err)
		}
		configPath = explicitPath
		return l.finalizeConfig(config, configPath)
	}

	// Otherwise, try default paths
	for _, path := range l.configPaths {
		if _, err := os.Stat(path); err == nil {
			config, err = LoadConfig(path)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %w", path, err)
			}
			configPath = path
			break
		}
	}

	// If no config file found, use defaults
	if config == nil {
		config = DefaultConfig()
		configPath = "."
	}

	return l.finalizeConfig(config, configPath)
}

// finalizeConfig performs final configuration steps
func (l *ConfigLoader) finalizeConfig(config *AppConfig, configPath string) (*AppConfig, error) {
	configDir := filepath.Dir(configPath)

	// Load fonts.toml if it exists and merge with config
	if err := l.loadFontsConfig(config, configDir); err != nil {
		return nil, fmt.Errorf("failed to load fonts config: %w", err)
	}

	return config, nil
}

// loadFontsConfig loads fonts from fonts.toml and merges with main config
func (l *ConfigLoader) loadFontsConfig(config *AppConfig, baseDir string) error {
	// Fonts are loaded directly by the font registry from fonts.toml
	// This method is kept for potential future use
	return nil
}

// SaveConfig saves the configuration to a TOML file
func SaveConfig(config *AppConfig, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Open file for writing
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close config file: %v\n", err)
		}
	}()

	// Encode to TOML
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// GenerateDefaultConfig generates a default configuration file
func GenerateDefaultConfig(path string) error {
	config := DefaultConfig()
	return SaveConfig(config, path)
}
