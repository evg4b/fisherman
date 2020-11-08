package utils

import (
	"os/exec"
	"path/filepath"
)

func NormalizePath(binary string) string {
	base := filepath.Base(binary)
	path, err := exec.LookPath(base)
	if err != nil || IsEmpty(path) {
		return binary
	}

	return base
}
