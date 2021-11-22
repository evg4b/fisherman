package utils

import "fmt"

// MergeEnvs concatenates environment variables list with passed map[string]string.
func MergeEnvs(env []string, newVars map[string]string) []string {
	envList := []string{}
	envList = append(env, envList...)

	for key, value := range newVars {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	return envList
}
