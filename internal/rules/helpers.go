package rules

import (
	"fmt"
	"strings"
)

func plainErrorFormatter(e []error) string {
	parts := []string{}
	for _, err := range e {
		parts = append(parts, err.Error())
	}

	return strings.Join(parts, "\n")
}

const maxPrefixLength = 14

func normalizePrefix(prefix string) string {
	if len(prefix) > maxPrefixLength {
		return fmt.Sprintf("%s...", prefix[:maxPrefixLength])
	}

	return prefix
}
