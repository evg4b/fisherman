package utils

import (
	"strings"
)

func Filter(slice []string, predicate func(string) bool) []string {
	resultSlice := make([]string, 0, len(slice))

	for i := range slice {
		if predicate(slice[i]) {
			resultSlice = append(resultSlice, slice[i])
		}
	}

	return resultSlice
}

func Contains(collection []string, value string) bool {
	for _, item := range collection {
		if strings.EqualFold(item, value) {
			return true
		}
	}

	return false
}
