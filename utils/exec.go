package utils

import "os/exec"

// IsCommandExists returns true when command registered in path as global tool
func IsCommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)

	return err == nil
}