package path

import (
	"os"
	"path/filepath"
	"strings"
)

func IsRegisteredInPath(path, app string) (bool, error) {
	normalizedAppPath, err := filepath.Abs(filepath.Dir(app))
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

		if matched {
			return true, nil
		}
	}

	return false, nil
}
