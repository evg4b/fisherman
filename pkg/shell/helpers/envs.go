package helpers

import (
	"fmt"
	"strings"
)

const keyValuePartsCount = 2

// MergeEnv concatenates environment variables list with passed map[string]string.
func MergeEnv(env []string, vars map[string]string) []string {
	if len(vars) == 0 {
		return env
	}

	parsedEnv := map[string]string{}
	for _, variable := range env {
		parts := strings.SplitN(variable, "=", keyValuePartsCount)
		parsedEnv[parts[0]] = parts[1]
	}

	for key, value := range vars {
		parsedEnv[key] = value
	}

	resultEnv := []string{}
	for key, value := range parsedEnv {
		resultEnv = append(resultEnv, fmt.Sprintf("%s=%s", key, value))
	}

	return resultEnv
}
