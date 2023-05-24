package create

import (
	"regexp"
	"strings"
)

var sanitize_regexp *regexp.Regexp

func init() {
	sanitize_regexp = regexp.MustCompile(`\s+`)
}

// Sanitize takes a string parameter and returns a modified string by replacing
// multiple consecutive whitespace characters with a single space, and trimming
// leading and trailing whitespace.
func Sanitize(input string) string {
	// Replace multiple consecutive whitespace characters with a single space
	sanitized := sanitize_regexp.ReplaceAllString(input, " ")

	// Trim leading and trailing whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// IsSanitizeCol returns true if the specifed column name is found in
// the list of columns that need to be sanitized.
func IsSanitizeCol(colName string) bool {
	for _, san := range sanitizeCols {
		if colName == san {
			return true
		}
	}
	return false
}
