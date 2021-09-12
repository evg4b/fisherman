package rules

import "strings"

func plainErrorFormatter(e []error) string {
	parts := []string{}
	for _, err := range e {
		parts = append(parts, err.Error())
	}

	return strings.Join(parts, "\n")
}
