package testutils

import (
	"os"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
)

func FsFromMap(t *testing.T, files map[string]string) billy.Filesystem {
	t.Helper()

	fs := memfs.New()
	for path, content := range files {
		err := util.WriteFile(fs, path, []byte(content), os.ModePerm)
		if err != nil {
			t.Error(err)
		}
	}

	return fs
}

func FsFromSlice(t *testing.T, files []string) billy.Filesystem {
	t.Helper()

	fs := memfs.New()
	for _, path := range files {
		err := util.WriteFile(fs, path, []byte("test"), os.ModePerm)
		if err != nil {
			t.Error(err)
		}
	}

	return fs
}
