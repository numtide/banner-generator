package resources

import (
	"os"
	"path/filepath"
)

// ResourceLocator helps find resources relative to the binary
type ResourceLocator struct {
	searchPaths []string
}

// NewResourceLocator creates a new resource locator
func NewResourceLocator() *ResourceLocator {
	locator := &ResourceLocator{
		searchPaths: []string{
			".",         // Current directory
			"templates", // Templates subdirectory
			"fonts",     // Fonts subdirectory
		},
	}

	// Add paths relative to executable for deployed scenarios
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		// For binaries in bin/ directory
		if filepath.Base(exeDir) == "bin" {
			projectRoot := filepath.Dir(exeDir)
			locator.searchPaths = append(locator.searchPaths,
				filepath.Join(projectRoot, "templates"),
				filepath.Join(projectRoot, "fonts"),
			)
		}
	}

	return locator
}

// FindFile searches for a file in all search paths
func (rl *ResourceLocator) FindFile(filename string) string {
	// If absolute path, return as-is if it exists
	if filepath.IsAbs(filename) {
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}

	// Search in all paths
	for _, searchPath := range rl.searchPaths {
		fullPath := filepath.Join(searchPath, filename)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	// Return original filename as fallback
	return filename
}

// FindTemplate finds a template file
func (rl *ResourceLocator) FindTemplate(name string) string {
	// Try with common extensions if not provided
	if filepath.Ext(name) == "" {
		for _, ext := range []string{".svg.mustache", ".mustache", ".svg"} {
			if found := rl.FindFile(name + ext); found != name+ext {
				return found
			}
		}
	}
	return rl.FindFile(name)
}

// FindFont finds a font file
func (rl *ResourceLocator) FindFont(name string) string {
	// Try with common extensions if not provided
	if filepath.Ext(name) == "" {
		for _, ext := range []string{".ttf", ".otf", ".woff", ".woff2"} {
			if found := rl.FindFile(name + ext); found != name+ext {
				return found
			}
		}
	}
	return rl.FindFile(name)
}
