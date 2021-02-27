package utils

import (
	"os/exec"
	"path/filepath"
	"time"
)

func NormalizePath(binary string) (string, bool) {
	base := filepath.Base(binary)
	path, err := exec.LookPath(base)
	if err != nil || IsEmpty(path) {
		return binary, true
	}

	return base, false
}

func ExecWithTime(runFunction func() error) (time.Duration, error) {
	start := time.Now()
	err := runFunction()

	return time.Since(start), err
}
