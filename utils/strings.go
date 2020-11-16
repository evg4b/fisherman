package utils

import "strings"

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func IsNotEmpty(value string) bool {
	return !IsEmpty(value)
}

func OriginalOrNA(path string) string {
	if IsEmpty(path) {
		return "N/A"
	}

	return path
}

func GetOrDefault(value string, defaultValue string) string {
	if IsNotEmpty(value) {
		return value
	}

	return defaultValue
}
