package utils

import "strings"

// IsEmpty return when string is empty or witespace
func IsEmpty(value string) bool {
	return len(strings.Trim(value, " ")) == 0
}
