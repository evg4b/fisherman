package fs

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

type Accessor struct {
}

func (f *Accessor) Exist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func NewAccessor() *Accessor {
	return &Accessor{}
}

func (f *Accessor) Read(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (f *Accessor) Reader(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(file), nil
}

func (f *Accessor) Delete(path string) error {
	return os.Remove(path)
}

func (f *Accessor) Write(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0600)
}
