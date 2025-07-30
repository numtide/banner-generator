package fonts

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Font represents a font with its various format files
type Font struct {
	Name     string            // Font name (e.g., "GT Pressura")
	Family   string            // CSS font-family name
	Variants map[string]string // Map of format to file path (e.g., "woff2" -> "path/to/font.woff2")
}

// GetFontPath returns the path to the font file, preferring WOFF format
func (f *Font) GetFontPath() string {
	// Try formats in order of preference
	formats := []string{"woff", "woff2", "ttf", "otf"}
	for _, format := range formats {
		if path, ok := f.Variants[format]; ok {
			return path
		}
	}
	return ""
}

// Registry manages fonts for both web serving and local access
type Registry struct {
	fonts   map[string]*Font  // Key is font family name
	aliases map[string]string // Font name aliases (e.g., "GT Pressura" -> "gt-pressura")
	baseDir string            // Base directory for font files
}

// NewRegistry creates a new font registry
func NewRegistry(baseDir string) *Registry {
	return &Registry{
		fonts:   make(map[string]*Font),
		aliases: make(map[string]string),
		baseDir: baseDir,
	}
}

// RegisterFont adds a font to the registry
func (r *Registry) RegisterFont(family string, font *Font) {
	font.Family = family
	r.fonts[family] = font
}

// RegisterAlias adds an alias for a font family
func (r *Registry) RegisterAlias(alias, family string) {
	r.aliases[alias] = family
}

// GetFont returns a font by family name or alias
func (r *Registry) GetFont(name string) *Font {
	// Try direct lookup
	if font, ok := r.fonts[name]; ok {
		return font
	}

	// Try alias lookup
	if family, ok := r.aliases[name]; ok {
		return r.fonts[family]
	}

	return nil
}

// GetFontPath returns the absolute path for a font format
func (r *Registry) GetFontPath(family, format string) (string, error) {
	font := r.GetFont(family)
	if font == nil {
		return "", fmt.Errorf("font family '%s' not found", family)
	}

	path, ok := font.Variants[format]
	if !ok {
		return "", fmt.Errorf("format '%s' not found for font '%s'", format, family)
	}

	// Return absolute path
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Join(r.baseDir, path), nil
}

// ServeHTTP implements http.Handler for serving font files
func (r *Registry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Extract path components
	path := strings.TrimPrefix(req.URL.Path, "/fonts/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.NotFound(w, req)
		return
	}

	// Expected format: /fonts/{family}/{filename}
	family := parts[0]
	filename := parts[1]

	// Determine format from filename
	format := ""
	switch {
	case strings.HasSuffix(filename, ".woff2"):
		format = "woff2"
	case strings.HasSuffix(filename, ".woff"):
		format = "woff"
	case strings.HasSuffix(filename, ".ttf"):
		format = "ttf"
	case strings.HasSuffix(filename, ".otf"):
		format = "otf"
	case strings.HasSuffix(filename, ".eot"):
		format = "eot"
	default:
		http.NotFound(w, req)
		return
	}

	// Get font path
	fontPath, err := r.GetFontPath(family, format)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	// Set appropriate content type
	contentType := getContentType(format)
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Set cache headers
	w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Serve the file
	http.ServeFile(w, req, fontPath)
}

// GetCSS generates @font-face CSS for a font family
func (r *Registry) GetCSS(family, baseURL string) (string, error) {
	font := r.GetFont(family)
	if font == nil {
		return "", fmt.Errorf("font family '%s' not found", family)
	}

	var css strings.Builder
	css.WriteString("@font-face {\n")
	css.WriteString(fmt.Sprintf("  font-family: '%s';\n", font.Family))
	css.WriteString("  src: ")

	// Build src with format hints
	var sources []string

	// Preferred order: woff2, woff, ttf, otf
	formats := []string{"woff2", "woff", "ttf", "otf"}
	for _, format := range formats {
		if _, ok := font.Variants[format]; ok {
			url := fmt.Sprintf("%s/fonts/%s/%s.%s", baseURL, family, family, format)
			var formatHint string
			switch format {
			case "ttf":
				formatHint = "truetype"
			case "otf":
				formatHint = "opentype"
			default:
				formatHint = format
			}
			sources = append(sources, fmt.Sprintf("url('%s') format('%s')", url, formatHint))
		}
	}

	css.WriteString(strings.Join(sources, ",\n       "))
	css.WriteString(";\n")
	css.WriteString("  font-weight: normal;\n")
	css.WriteString("  font-style: normal;\n")
	css.WriteString("}\n")

	return css.String(), nil
}

// LoadFontData loads font file data for embedding
func (r *Registry) LoadFontData(family, format string) ([]byte, error) {
	fontPath, err := r.GetFontPath(family, format)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(fontPath)
}

func getContentType(format string) string {
	switch format {
	case "woff2":
		return "font/woff2"
	case "woff":
		return "font/woff"
	case "ttf":
		return "font/ttf"
	case "otf":
		return "font/otf"
	case "eot":
		return "application/vnd.ms-fontobject"
	default:
		return "application/octet-stream"
	}
}

// DefaultRegistry creates a registry with the default GT Pressura font
func DefaultRegistry(baseDir string) *Registry {
	registry := NewRegistry(baseDir)

	// Register GT Pressura
	registry.RegisterFont("gt-pressura", &Font{
		Name: "GT Pressura Regular",
		Variants: map[string]string{
			"ttf":   "gt-pressura-regular.ttf",
			"woff":  "web/gt-pressura-regular.woff",
			"woff2": "web/gt-pressura-regular.woff2",
		},
	})

	// Register common aliases
	registry.RegisterAlias("GT Pressura", "gt-pressura")
	registry.RegisterAlias("gt pressura", "gt-pressura")
	registry.RegisterAlias("GTpressura", "gt-pressura")

	return registry
}
