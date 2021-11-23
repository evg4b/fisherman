package utils

import "fmt"

// MergeEnv concatenates environment variables list with passed map[string]string.
func MergeEnv(env []string, vars map[string]string) []string {
	newEnv := []string{}
	newEnv = append(env, newEnv...)

	for key, value := range vars {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", key, value))
	}

	return newEnv
}
