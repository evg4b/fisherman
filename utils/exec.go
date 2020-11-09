package utils

import (
	"os/exec"
	"path/filepath"
	"time"
)

func NormalizePath(binary string) string {
	base := filepath.Base(binary)
	path, err := exec.LookPath(base)
	if err != nil || IsEmpty(path) {
		return binary
	}

	return base
}

func ExecWithTime(runFunction func() error) (time.Duration, error) {
	start := time.Now()
	err := runFunction()

	if err != nil {
		return time.Duration(0), err
	}

	return time.Since(start), err
}
