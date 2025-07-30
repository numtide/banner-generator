package utils

import "fmt"

// FormatCount formats a number with K suffix for thousands
func FormatCount(count int) string {
	if count >= 1000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}
