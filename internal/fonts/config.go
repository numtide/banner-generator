package fonts

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the font configuration from TOML file
type Config struct {
	Fonts []FontConfig `toml:"fonts"`
}

// FontConfig represents a single font configuration
type FontConfig struct {
	Family   string            `toml:"family"`
	Name     string            `toml:"name"`
	Aliases  []string          `toml:"aliases"`
	Variants map[string]string `toml:"variants"`
}

// LoadConfig loads font configuration from a TOML file
func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to load font config: %w", err)
	}
	return &config, nil
}

// BuildRegistry creates a font registry from configuration
func BuildRegistry(config *Config, baseDir string) *Registry {
	registry := NewRegistry(baseDir)

	for _, fontConfig := range config.Fonts {
		// Register the font
		font := &Font{
			Name:     fontConfig.Name,
			Family:   fontConfig.Family,
			Variants: fontConfig.Variants,
		}
		registry.RegisterFont(fontConfig.Family, font)

		// Register aliases
		for _, alias := range fontConfig.Aliases {
			registry.RegisterAlias(alias, fontConfig.Family)
		}
	}

	return registry
}

// LoadRegistryFromTOML loads a font registry from a TOML configuration file
func LoadRegistryFromTOML(configPath, baseDir string) (*Registry, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return BuildRegistry(config, baseDir), nil
}

// DefaultRegistryWithConfig creates a registry from fonts.toml or falls back to defaults
func DefaultRegistryWithConfig(baseDir string) *Registry {
	configPath := filepath.Join(baseDir, "fonts.toml")

	// Try to load from config file
	if registry, err := LoadRegistryFromTOML(configPath, baseDir); err == nil {
		return registry
	}

	// Fall back to default registry
	return DefaultRegistry(baseDir)
}
