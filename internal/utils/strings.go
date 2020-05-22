package utils

import (
	"errors"
	"strings"
)

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func FirstNotEmpty(values ...string) string {
	if len(values) == 0 {
		panic(errors.New("sequence contains no elements"))
	}

	for _, value := range values {
		if !IsEmpty(value) {
			return value
		}
	}

	return values[len(values)-1]
}
