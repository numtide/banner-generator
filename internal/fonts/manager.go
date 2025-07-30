package fonts

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Manager handles font operations for banner generation
type Manager interface {
	GetFont(family string) *Font
	GetFontData(fontPath string) (string, error)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// DefaultManager implements Manager using a Registry
type DefaultManager struct {
	registry *Registry
	baseDir  string
}

// NewManager creates a new font manager
func NewManager(fontDir string) Manager {
	registry := DefaultRegistryWithConfig(fontDir)
	return &DefaultManager{
		registry: registry,
		baseDir:  fontDir,
	}
}

// GetFont returns a font by family name or alias
func (m *DefaultManager) GetFont(family string) *Font {
	return m.registry.GetFont(family)
}

// GetFontData returns base64-encoded font data
func (m *DefaultManager) GetFontData(fontPath string) (string, error) {
	fullPath := fontPath
	if !filepath.IsAbs(fontPath) {
		fullPath = filepath.Join(m.baseDir, fontPath)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read font file: %w", err)
	}

	// Determine MIME type
	mimeType := "font/ttf"
	switch strings.ToLower(filepath.Ext(fontPath)) {
	case ".woff":
		mimeType = "font/woff"
	case ".woff2":
		mimeType = "font/woff2"
	case ".otf":
		mimeType = "font/otf"
	}

	// Return as data URI
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

// ServeHTTP implements http.Handler for serving font files
func (m *DefaultManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.registry.ServeHTTP(w, r)
}
