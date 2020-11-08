package filesystem

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
)

type LocalFileSystem struct {
}

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

func (f *LocalFileSystem) Reader(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(file), nil
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

func (f *LocalFileSystem) Chown(path string, user *user.User) error {
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return err
	}

	return os.Chown(path, uid, gid)
}
