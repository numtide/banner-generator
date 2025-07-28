package fonts

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/numtide/banner-generator/internal/utils"
)

// FontResolver extracts font requirements from SVG and resolves them
type FontResolver struct {
	registry *Registry
}

// NewFontResolver creates a new font resolver
func NewFontResolver(registry *Registry) *FontResolver {
	return &FontResolver{
		registry: registry,
	}
}

// FontRequirement represents a required font from the SVG
type FontRequirement struct {
	Family string
	Weight string // normal, bold, etc.
	Style  string // normal, italic, etc.
}

// ResolvedFonts contains the resolved font data
type ResolvedFonts struct {
	Requirements []FontRequirement
	EmbeddedCSS  string            // CSS with embedded base64 fonts
	WebCSS       string            // CSS with web font URLs
	FontData     map[string][]byte // Font family -> font data for embedding
}

// ExtractFontsFromSVG parses SVG and extracts all font-family references
func (fr *FontResolver) ExtractFontsFromSVG(svgContent string) ([]FontRequirement, error) {
	requirements := make(map[string]FontRequirement)

	// Parse SVG to find font references
	decoder := xml.NewDecoder(strings.NewReader(svgContent))

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing SVG: %w", err)
		}

		switch elem := token.(type) {
		case xml.StartElement:
			// Check text and tspan elements
			if elem.Name.Local == "text" || elem.Name.Local == "tspan" {
				fontFamily := ""
				fontWeight := "normal"
				fontStyle := "normal"

				// Check attributes
				for _, attr := range elem.Attr {
					switch attr.Name.Local {
					case "font-family":
						fontFamily = cleanFontFamily(attr.Value)
					case "font-weight":
						fontWeight = attr.Value
					case "font-style":
						fontStyle = attr.Value
					}
				}

				// Check style attribute
				for _, attr := range elem.Attr {
					if attr.Name.Local == "style" {
						styleProps := parseStyleAttribute(attr.Value)
						if family, ok := styleProps["font-family"]; ok {
							fontFamily = cleanFontFamily(family)
						}
						if weight, ok := styleProps["font-weight"]; ok {
							fontWeight = weight
						}
						if style, ok := styleProps["font-style"]; ok {
							fontStyle = style
						}
					}
				}

				if fontFamily != "" {
					key := fmt.Sprintf("%s-%s-%s", fontFamily, fontWeight, fontStyle)
					requirements[key] = FontRequirement{
						Family: fontFamily,
						Weight: fontWeight,
						Style:  fontStyle,
					}
				}
			}
		}
	}

	// Convert map to slice
	var result []FontRequirement
	for _, req := range requirements {
		result = append(result, req)
	}

	return result, nil
}

// ResolveFonts takes font requirements and resolves them to actual font files
func (fr *FontResolver) ResolveFonts(requirements []FontRequirement, baseURL string) (*ResolvedFonts, error) {
	resolved := &ResolvedFonts{
		Requirements: requirements,
		FontData:     make(map[string][]byte),
	}

	var embeddedStyles []string
	var webStyles []string

	for _, req := range requirements {
		// Try to find font in registry
		font := fr.registry.GetFont(utils.NormalizeFontName(req.Family))
		if font == nil {
			// Try with original name
			font = fr.registry.GetFont(req.Family)
		}

		if font == nil {
			return nil, fmt.Errorf("font not found: %s", req.Family)
		}

		// For embedded fonts
		if data, err := fr.registry.LoadFontData(font.Family, "ttf"); err == nil {
			resolved.FontData[req.Family] = data
			embeddedStyles = append(embeddedStyles, fmt.Sprintf(`@font-face {
  font-family: '%s';
  src: url('data:font/truetype;base64,%s') format('truetype');
  font-weight: %s;
  font-style: %s;
}`, req.Family, "{{"+utils.NormalizeFontName(req.Family)+"FontData}}", req.Weight, req.Style))
		}

		// For web fonts
		if baseURL != "" {
			css, err := fr.registry.GetCSS(font.Family, baseURL)
			if err == nil {
				webStyles = append(webStyles, css)
			}
		}
	}

	resolved.EmbeddedCSS = strings.Join(embeddedStyles, "\n")
	resolved.WebCSS = strings.Join(webStyles, "\n")

	return resolved, nil
}

// cleanFontFamily removes quotes and extra spaces from font family names
func cleanFontFamily(family string) string {
	family = strings.TrimSpace(family)
	family = strings.Trim(family, "'\"")
	return family
}

// parseStyleAttribute parses CSS style attribute into key-value pairs
func parseStyleAttribute(style string) map[string]string {
	props := make(map[string]string)

	declarations := strings.Split(style, ";")
	for _, decl := range declarations {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}

		parts := strings.SplitN(decl, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			props[key] = value
		}
	}

	return props
}
