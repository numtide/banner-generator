package banner

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/numtide/banner-generator/internal/fonts"
	"github.com/numtide/banner-generator/internal/utils"
)

// AutoFontBuilder creates SVG banners with automatic font resolution
type AutoFontBuilder struct {
	registry     *fonts.Registry
	resolver     *fonts.FontResolver
	templatePath string
	baseURL      string // Base URL for web fonts (optional)
}

// NewAutoFontBuilder creates a new builder with automatic font resolution
func NewAutoFontBuilder(registry *fonts.Registry, templatePath string, baseURL string) *AutoFontBuilder {
	return &AutoFontBuilder{
		registry:     registry,
		resolver:     fonts.NewFontResolver(registry),
		templatePath: templatePath,
		baseURL:      baseURL,
	}
}

// BuildSVG creates an SVG banner with the provided data
func (b *AutoFontBuilder) BuildSVG(data *BannerData) (string, error) {
	// First, render the template with basic data
	templateData := b.prepareTemplateData(data)

	// Initial render to get SVG content
	initialSVG, err := mustache.RenderFile(b.templatePath, templateData)
	if err != nil {
		return "", fmt.Errorf("failed to render initial template: %w", err)
	}

	// Extract font requirements from the SVG
	requirements, err := b.resolver.ExtractFontsFromSVG(initialSVG)
	if err != nil {
		return "", fmt.Errorf("failed to extract fonts: %w", err)
	}

	// Resolve fonts
	resolved, err := b.resolver.ResolveFonts(requirements, b.baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to resolve fonts: %w", err)
	}

	// Add font data to template
	if b.baseURL != "" {
		// Web fonts mode
		templateData["useWebFonts"] = true
		templateData["fontCSS"] = resolved.WebCSS
	} else {
		// Embedded fonts mode
		templateData["useEmbeddedFont"] = true

		// Add base64 encoded font data for each font
		for family, data := range resolved.FontData {
			key := utils.NormalizeFontName(family) + "FontData"
			templateData[key] = base64.StdEncoding.EncodeToString(data)
		}

		// Build embedded CSS with template variables
		var embeddedStyles []string
		for _, req := range resolved.Requirements {
			normalizedName := utils.NormalizeFontName(req.Family)
			embeddedStyles = append(embeddedStyles, fmt.Sprintf(`@font-face {
  font-family: '%s';
  src: url('data:font/truetype;base64,{{%sFontData}}') format('truetype');
  font-weight: %s;
  font-style: %s;
}`, req.Family, normalizedName, req.Weight, req.Style))
		}
		templateData["fontCSS"] = strings.Join(embeddedStyles, "\n")
	}

	// Final render with font data
	finalSVG, err := mustache.RenderFile(b.templatePath, templateData)
	if err != nil {
		return "", fmt.Errorf("failed to render final template: %w", err)
	}

	return finalSVG, nil
}

// prepareTemplateData prepares the basic template data
func (b *AutoFontBuilder) prepareTemplateData(data *BannerData) map[string]interface{} {
	repoName := data.RepoName

	// Calculate font size based on repository name (not full path)
	repoNameFontSize := 153
	if len(data.RepoName) > 12 {
		repoNameFontSize = 120
	}
	if len(data.RepoName) > 20 {
		repoNameFontSize = 100
	}
	if len(data.RepoName) > 30 {
		repoNameFontSize = 80
	}

	// Prepare template data
	templateData := map[string]interface{}{
		"repoName":         repoName,
		"repoNameFontSize": repoNameFontSize,
		"hasDescription":   data.RepoDescription != "",
		"hasLanguage":      data.Language != "",
		"stars":            formatCount(data.Stars),
		"forks":            formatCount(data.Forks),
		"hasStats":         data.Stars > 0 || data.Forks > 0,
		"owner":            data.Owner,
	}

	// Add description if available
	if data.RepoDescription != "" {
		// Split description into lines (max 50 chars per line)
		descLines := splitDescription(data.RepoDescription, 50)

		// Create line data with Y positions
		var lines []map[string]interface{}
		startY := 466.802
		lineHeight := 63.0

		for i, line := range descLines {
			if i > 1 { // Max 2 lines
				break
			}
			lines = append(lines, map[string]interface{}{
				"text": line,
				"y":    fmt.Sprintf("%.3f", startY+float64(i)*lineHeight),
			})
		}

		templateData["descriptionLines"] = lines
		templateData["description"] = data.RepoDescription // Keep original for compatibility
	}

	// Add language if available
	if data.Language != "" {
		templateData["language"] = data.Language
	}

	return templateData
}

// formatCount formats a number for display (e.g., 1234 -> "1.2k")
func formatCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	} else if count < 10000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	} else if count < 1000000 {
		return fmt.Sprintf("%dk", count/1000)
	} else {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	}
}

// splitDescription splits a description into lines that fit within maxLength
func splitDescription(desc string, maxLength int) []string {
	words := strings.Fields(desc)
	var lines []string
	var currentLine []string
	currentLength := 0

	for _, word := range words {
		wordLen := len(word)
		if currentLength > 0 && currentLength+1+wordLen > maxLength {
			// Start new line
			lines = append(lines, strings.Join(currentLine, " "))
			currentLine = []string{word}
			currentLength = wordLen
		} else {
			// Add to current line
			if currentLength > 0 {
				currentLength++ // Space
			}
			currentLine = append(currentLine, word)
			currentLength += wordLen
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	return lines
}
