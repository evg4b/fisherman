package utils

import (
	"io/ioutil"
	"os"

	"github.com/go-errors/errors"
	"github.com/go-git/go-billy/v5"
)

func ReadFileAsString(fs billy.Filesystem, filepath string) (string, error) {
	file, err := fs.Open(filepath)
	if err != nil {
		return "", errors.Errorf("open %s: %w", filepath, err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

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
