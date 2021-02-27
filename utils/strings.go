package utils

import "strings"

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func OriginalOrNA(path string) string {
	if IsEmpty(path) {
		return "N/A"
	}

	return path
}

func GetOrDefault(value string, defaultValue string) string {
	if !IsEmpty(value) {
		return value
	}

	return defaultValue
}

func Contains(collection []string, value string) bool {
	for _, item := range collection {
		if strings.EqualFold(item, value) {
			return true
		}
	}

	return false
}
