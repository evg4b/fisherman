package utils

func StringMapToInterfaceMap(src map[string]string) map[string]any {
	dest := map[string]any{}
	for key, value := range src {
		dest[key] = value
	}

	return dest
}
