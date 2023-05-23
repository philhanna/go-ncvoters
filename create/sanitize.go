package create

import (
	"regexp"
	"strings"
)

// Sanitize takes a string parameter and returns a modified string by replacing
// multiple consecutive whitespace characters with a single space, and trimming
// leading and trailing whitespace.
func Sanitize(input string) string {
	// Replace multiple consecutive whitespace characters with a single space
	regex := regexp.MustCompile(`\s+`)
	sanitized := regex.ReplaceAllString(input, " ")

	// Trim leading and trailing whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}
