package path

import (
	"os"
	"path/filepath"
	"strings"
)

// IsRegisteredInPath checks if a executable is registered in the path variable
func IsRegisteredInPath(path, executablePath string) (bool, error) {
	normalizedAppPath, err := filepath.Abs(filepath.Dir(executablePath))
	if err != nil {
		return false, err
	}
	parts := strings.Split(path, string(os.PathListSeparator))
	for _, pathItem := range parts {
		normalized, err := filepath.Abs(pathItem)
		if err != nil {
			return false, nil
		}
		matched, err := filepath.Match(normalized, normalizedAppPath)
		if err != nil {
			return false, err
		}
		return matched, nil
	}
	return false, nil
}
