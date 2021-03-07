package utils

import (
	"github.com/spf13/afero"
)

func ReadFileAsString(fs afero.Fs, filepath string) (string, error) {
	data, err := afero.ReadFile(fs, filepath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
