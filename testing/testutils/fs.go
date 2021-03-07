package testutils

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func FsFromMap(t *testing.T, files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()
	for path, content := range files {
		err := afero.WriteFile(fs, path, []byte(content), os.ModePerm)
		if err != nil {
			t.Error(err)
		}
	}

	return fs
}

func FsFromSlice(t *testing.T, files []string) afero.Fs {
	fs := afero.NewMemMapFs()
	for _, path := range files {
		err := afero.WriteFile(fs, path, []byte("test"), os.ModePerm)
		if err != nil {
			t.Error(err)
		}
	}

	return fs
}
