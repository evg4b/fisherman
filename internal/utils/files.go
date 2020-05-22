package utils

import (
	"os"

	"github.com/go-git/go-billy/v5"
)

func Exists(fs billy.Filesystem, filepath string) (bool, error) {
	_, err := fs.Stat(filepath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
