package utils

import "strings"

// NormalizeFontName converts font names to safe identifiers
// Example: "GT Pressura" -> "gt_pressura"
func NormalizeFontName(name string) string {
	// Replace spaces and special chars with underscores
	normalized := strings.ReplaceAll(name, " ", "_")
	normalized = strings.ReplaceAll(normalized, "-", "_")
	return strings.ToLower(normalized)
}
