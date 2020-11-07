package utils

// Filter filters string slice by predicate
func Filter(slice []string, predicate func(string) bool) []string {
	resultSlice := []string{}

	for i := range slice {
		if predicate(slice[i]) {
			resultSlice = append(resultSlice, slice[i])
		}
	}

	return resultSlice
}
