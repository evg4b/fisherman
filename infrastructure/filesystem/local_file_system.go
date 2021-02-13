package filesystem

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type LocalFileSystem struct{}

func NewLocalFileSystem() *LocalFileSystem {
	return &LocalFileSystem{}
}

func (f *LocalFileSystem) Exist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func (f *LocalFileSystem) Read(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (f *LocalFileSystem) Reader(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (f *LocalFileSystem) Delete(path string) error {
	return os.Remove(path)
}

func (f *LocalFileSystem) Write(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0600)
}

func (f *LocalFileSystem) Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

func (f *LocalFileSystem) Find(folder string, globs []string) ([]string, error) {
	files := []string{}

	err := filepath.Walk(folder, func(file string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}

		for _, glob := range globs {
			matched, err := path.Match(glob, info.Name())
			if err != nil {
				return err
			}

			if matched {
				files = append(files, file)

				return nil
			}
		}

		return nil
	})

	return files, err
}
