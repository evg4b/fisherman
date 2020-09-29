package utils

import "strings"

// IsEmpty returns when string is empty or witespace
func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

// IsNotEmpty returns when string is empty or witespace
func IsNotEmpty(value string) bool {
	return !IsEmpty(value)
}

// OriginalOrNA returns original string when string is not empty. Else returns "N/A".
func OriginalOrNA(path string) string {
	if IsEmpty(path) {
		return "N/A"
	}

	return path
}
