package utils

func StringMapToInterfaceMap(src map[string]string) map[string]interface{} {
	dest := map[string]interface{}{}
	for key, value := range src {
		dest[key] = value
	}

	return dest
}
