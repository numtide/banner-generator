package svg

import (
	"strings"
	"testing"
)

func TestUpdateTextByID(t *testing.T) {
	svg := `<svg>
		<text id="title" x="10" y="20">Old Title</text>
		<text id="subtitle">Old Subtitle</text>
	</svg>`

	doc := NewSimpleDocument(svg)

	// Update title
	err := doc.UpdateTextByID("title", "New Title")
	if err != nil {
		t.Errorf("Failed to update title: %v", err)
	}

	result := doc.String()
	if !strings.Contains(result, "New Title") {
		t.Error("Result doesn't contain new title")
	}
	if strings.Contains(result, "Old Title") {
		t.Error("Result still contains old title")
	}

	// Update subtitle
	err = doc.UpdateTextByID("subtitle", "New Subtitle")
	if err != nil {
		t.Errorf("Failed to update subtitle: %v", err)
	}

	result = doc.String()
	if !strings.Contains(result, "New Subtitle") {
		t.Error("Result doesn't contain new subtitle")
	}
}

func TestUpdateMultilineText(t *testing.T) {
	svg := `<svg>
		<text id="description" x="50" y="100">
			<tspan>Old line 1</tspan>
			<tspan>Old line 2</tspan>
		</text>
	</svg>`

	doc := NewSimpleDocument(svg)

	lines := []string{"New line 1", "New line 2", "New line 3"}
	err := doc.UpdateMultilineText("description", lines)
	if err != nil {
		t.Errorf("Failed to update multiline text: %v", err)
	}

	result := doc.String()
	for _, line := range lines {
		if !strings.Contains(result, line) {
			t.Errorf("Result doesn't contain line: %s", line)
		}
	}

	// Check tspan structure
	if !strings.Contains(result, `<tspan x="50" dy="0">New line 1</tspan>`) {
		t.Error("First tspan not formatted correctly")
	}
	if !strings.Contains(result, `<tspan x="50" dy="1.2em">New line 2</tspan>`) {
		t.Error("Second tspan not formatted correctly")
	}
}

func TestHideElementByID(t *testing.T) {
	svg := `<svg>
		<g id="stats-group">
			<text>Some stats</text>
		</g>
	</svg>`

	doc := NewSimpleDocument(svg)

	err := doc.HideElementByID("stats-group")
	if err != nil {
		t.Errorf("Failed to hide element: %v", err)
	}

	result := doc.String()
	if !strings.Contains(result, `visibility="hidden"`) {
		t.Error("Element not hidden")
	}
}

func TestInjectCSS(t *testing.T) {
	svg := `<svg>
		<style id="font-css">
			/* Existing CSS */
		</style>
	</svg>`

	doc := NewSimpleDocument(svg)

	newCSS := `@font-face { font-family: 'Test'; }`
	err := doc.InjectCSS("font-css", newCSS)
	if err != nil {
		t.Errorf("Failed to inject CSS: %v", err)
	}

	result := doc.String()
	if !strings.Contains(result, newCSS) {
		t.Error("CSS not injected")
	}
	if !strings.Contains(result, "/* Existing CSS */") {
		t.Error("Existing CSS removed")
	}
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello & World", "Hello &amp; World"},
		{"<tag>", "&lt;tag&gt;"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"'single'", "&apos;single&apos;"},
		{"Normal text", "Normal text"},
	}

	for _, tt := range tests {
		result := escapeXML(tt.input)
		if result != tt.expected {
			t.Errorf("escapeXML(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
