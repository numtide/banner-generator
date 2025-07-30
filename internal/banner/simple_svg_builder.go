package banner

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/numtide/banner-generator/internal/fonts"
	"github.com/numtide/banner-generator/internal/github"
	"github.com/numtide/banner-generator/internal/svg"
	"github.com/numtide/banner-generator/internal/utils"
)

// SimpleSVGBuilder builds banners using simple string manipulation
type SimpleSVGBuilder struct {
	fontManager     fonts.Manager
	templatePath    string
	enableWebFonts  bool
	webFontsBaseURL string
}

// NewSimpleSVGBuilder creates a new simple SVG-based banner builder
func NewSimpleSVGBuilder(fontManager fonts.Manager, templatePath string, enableWebFonts bool, webFontsBaseURL string) *SimpleSVGBuilder {
	return &SimpleSVGBuilder{
		fontManager:     fontManager,
		templatePath:    templatePath,
		enableWebFonts:  enableWebFonts,
		webFontsBaseURL: webFontsBaseURL,
	}
}

// BuildBanner generates a banner for the given repository
func (b *SimpleSVGBuilder) BuildBanner(repo *github.Repository) (string, error) {
	// Load template
	templateContent, err := os.ReadFile(b.templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Create simple document
	doc := svg.NewSimpleDocument(string(templateContent))

	// Update repository name
	if err := doc.UpdateTextByID("repo-name", repo.Name); err != nil {
		// Ignore error - element might not exist in template
		_ = err
	}

	// Update description with multi-line support
	if repo.Description != "" {
		lines := wrapText(repo.Description, 50) // 50 chars per line
		if err := doc.UpdateMultilineText("description", lines); err != nil {
			// Ignore error - element might not exist in template
			_ = err
		}
	} else {
		// Hide description if empty
		_ = doc.HideElementByID("description")
	}

	// Update stats
	if repo.StargazersCount > 0 || repo.ForksCount > 0 {
		// Format star count
		stars := utils.FormatCount(repo.StargazersCount)
		_ = doc.UpdateTextByID("stats-stars", fmt.Sprintf("⭐ %s", stars))

		// Format fork count
		forks := utils.FormatCount(repo.ForksCount)
		_ = doc.UpdateTextByID("stats-forks", fmt.Sprintf("🍴 %s", forks))

		// Update language if available
		if repo.Language != "" {
			_ = doc.UpdateTextByID("stats-language", repo.Language)
		} else {
			_ = doc.HideElementByID("stats-language")
		}
	} else {
		// Hide stats group if no data
		_ = doc.HideElementByID("stats-group")
	}

	// Generate font CSS
	fontCSS, err := b.generateFontCSS(doc.String())
	if err != nil {
		return "", fmt.Errorf("failed to generate font CSS: %w", err)
	}

	// Inject font CSS
	if err := doc.InjectCSS("font-css", fontCSS); err != nil {
		// Try injecting into any style element
		_ = err
	}

	return doc.String(), nil
}

// generateFontCSS generates @font-face CSS for fonts used in the SVG
func (b *SimpleSVGBuilder) generateFontCSS(svgContent string) (string, error) {
	// Find all font-family references
	fontFamilies := make(map[string]bool)

	// Pattern to find font-family in attributes
	attrPattern := `font-family="([^"]*)"`
	attrRe := regexp.MustCompile(attrPattern)
	for _, match := range attrRe.FindAllStringSubmatch(svgContent, -1) {
		if len(match) > 1 {
			fontFamilies[match[1]] = true
		}
	}

	// Pattern to find font-family in style
	stylePattern := `font-family:\s*'?([^'";]+)'?`
	styleRe := regexp.MustCompile(stylePattern)
	for _, match := range styleRe.FindAllStringSubmatch(svgContent, -1) {
		if len(match) > 1 {
			family := strings.TrimSpace(match[1])
			family = strings.Trim(family, `"'`)
			fontFamilies[family] = true
		}
	}

	// Generate CSS for each font
	var cssBuilder strings.Builder
	for family := range fontFamilies {
		font := b.fontManager.GetFont(family)
		if font == nil {
			continue // Skip unknown fonts
		}

		if b.enableWebFonts {
			// Use web fonts with URLs
			cssBuilder.WriteString(b.generateWebFontCSS(font))
		} else {
			// Embed font data
			cssBuilder.WriteString(b.generateEmbeddedFontCSS(font))
		}
		cssBuilder.WriteString("\n")
	}

	return cssBuilder.String(), nil
}

// generateWebFontCSS generates @font-face CSS with URLs
func (b *SimpleSVGBuilder) generateWebFontCSS(font *fonts.Font) string {
	var css strings.Builder
	css.WriteString("@font-face {\n")
	css.WriteString(fmt.Sprintf("  font-family: '%s';\n", font.Family))
	css.WriteString("  src: ")

	var sources []string

	// WOFF2
	if woff2, ok := font.Variants["woff2"]; ok {
		url := b.webFontsBaseURL + "/fonts/" + filepath.Base(woff2)
		sources = append(sources, fmt.Sprintf("url('%s') format('woff2')", url))
	}

	// WOFF
	if woff, ok := font.Variants["woff"]; ok {
		url := b.webFontsBaseURL + "/fonts/" + filepath.Base(woff)
		sources = append(sources, fmt.Sprintf("url('%s') format('woff')", url))
	}

	// TTF
	if ttf, ok := font.Variants["ttf"]; ok {
		url := b.webFontsBaseURL + "/fonts/" + filepath.Base(ttf)
		sources = append(sources, fmt.Sprintf("url('%s') format('truetype')", url))
	}

	css.WriteString(strings.Join(sources, ",\n       "))
	css.WriteString(";\n")
	css.WriteString("  font-weight: normal;\n")
	css.WriteString("  font-style: normal;\n")
	css.WriteString("}\n")

	return css.String()
}

// generateEmbeddedFontCSS generates @font-face CSS with embedded data
func (b *SimpleSVGBuilder) generateEmbeddedFontCSS(font *fonts.Font) string {
	fontPath := font.GetFontPath()
	if fontPath == "" {
		return ""
	}

	fontData, err := b.fontManager.GetFontData(fontPath)
	if err != nil {
		return ""
	}

	format := "truetype"
	if strings.HasSuffix(fontPath, ".woff") {
		format = "woff"
	} else if strings.HasSuffix(fontPath, ".woff2") {
		format = "woff2"
	}

	return fmt.Sprintf(`@font-face {
  font-family: '%s';
  src: url(%s) format('%s');
  font-weight: normal;
  font-style: normal;
}`, font.Family, fontData, format)
}

// wrapText breaks text into lines based on character count
func wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder

	for i, word := range words {
		if i > 0 {
			// Check if adding this word would exceed the limit
			if currentLine.Len()+1+len(word) > maxWidth {
				// Start a new line
				lines = append(lines, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
			} else {
				// Add to current line
				currentLine.WriteString(" ")
				currentLine.WriteString(word)
			}
		} else {
			// First word
			currentLine.WriteString(word)
		}
	}

	// Add the last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}
