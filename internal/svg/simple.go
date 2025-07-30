package svg

import (
	"fmt"
	"regexp"
	"strings"
)

// SimpleDocument represents a simple SVG document for string-based manipulation
type SimpleDocument struct {
	content string
}

// NewSimpleDocument creates a new simple SVG document from a string
func NewSimpleDocument(content string) *SimpleDocument {
	return &SimpleDocument{content: content}
}

// UpdateTextByID updates the text content of an element with the given ID
func (d *SimpleDocument) UpdateTextByID(id, newText string) error {
	// Pattern to match text element with specific ID and capture everything between opening and closing tags
	pattern := fmt.Sprintf(`(?s)(<text[^>]*\sid="%s"[^>]*>).*?(</text>)`, id)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	if !re.MatchString(d.content) {
		return fmt.Errorf("element with id '%s' not found", id)
	}

	// Replace with opening tag + tspan with new text + closing tag
	d.content = re.ReplaceAllString(d.content, "${1}<tspan>"+escapeXML(newText)+"</tspan>${2}")
	return nil
}

// UpdateMultilineText updates a text element with multiple tspan children
func (d *SimpleDocument) UpdateMultilineText(id string, lines []string) error {
	// Find the text element (with multiline flag)
	pattern := fmt.Sprintf(`(?s)<text[^>]*\sid="%s"[^>]*>.*?</text>`, id)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	match := re.FindString(d.content)
	if match == "" {
		return fmt.Errorf("text element with id '%s' not found", id)
	}

	// Extract attributes from the text element
	attrPattern := `<text([^>]*)>`
	attrRe := regexp.MustCompile(attrPattern)
	attrMatch := attrRe.FindStringSubmatch(match)
	if len(attrMatch) < 2 {
		return fmt.Errorf("failed to parse text element attributes")
	}

	// Extract x coordinate for tspans
	xPattern := `x="([^"]*)"`
	xRe := regexp.MustCompile(xPattern)
	xMatch := xRe.FindStringSubmatch(match)
	x := "0"
	if len(xMatch) >= 2 {
		x = xMatch[1]
	}

	// Build new text element with tspans
	var newElement strings.Builder
	newElement.WriteString(fmt.Sprintf(`<text%s>`, attrMatch[1]))

	for i, line := range lines {
		if i == 0 {
			newElement.WriteString(fmt.Sprintf(`<tspan x="%s" dy="0">%s</tspan>`, x, escapeXML(line)))
		} else {
			newElement.WriteString(fmt.Sprintf(`<tspan x="%s" dy="1.2em">%s</tspan>`, x, escapeXML(line)))
		}
	}

	newElement.WriteString(`</text>`)

	// Replace the old element with the new one
	d.content = strings.Replace(d.content, match, newElement.String(), 1)
	return nil
}

// HideElementByID hides an element by adding visibility="hidden"
func (d *SimpleDocument) HideElementByID(id string) error {
	// Pattern to match any element with the given ID
	pattern := fmt.Sprintf(`(<[^>]+\sid="%s"[^>]*)>`, id)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	if !re.MatchString(d.content) {
		return fmt.Errorf("element with id '%s' not found", id)
	}

	// Add visibility="hidden" attribute
	d.content = re.ReplaceAllString(d.content, `${1} visibility="hidden">`)
	return nil
}

// InjectCSS injects CSS content into a style element with the given ID
func (d *SimpleDocument) InjectCSS(styleID, css string) error {
	// Pattern to match style element with ID (with multiline flag)
	pattern := fmt.Sprintf(`(?s)(<style[^>]*\sid="%s"[^>]*>)(.*?)(</style>)`, styleID)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	if !re.MatchString(d.content) {
		// Try without ID requirement
		pattern = `(?s)(<style[^>]*>)(.*?)(</style>)`
		re = regexp.MustCompile(pattern)
		if !re.MatchString(d.content) {
			return fmt.Errorf("style element not found")
		}
	}

	// Replace with existing content + new CSS
	d.content = re.ReplaceAllString(d.content, "${1}${2}\n"+css+"${3}")
	return nil
}

// String returns the SVG content as a string
func (d *SimpleDocument) String() string {
	return d.content
}

// escapeXML escapes special XML characters
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
