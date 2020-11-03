package infrastructure

import "io"

type User struct {
	UserName string
	Email    string
}

type Repository interface {
	GetCurrentBranch() (string, error)
	GetUser() (User, error)
	GetLastTag() (string, error)
}

type FileWriter interface {
	Write(path, content string) error
}

type FileReader interface {
	Read(path string) (string, error)
	Reader(path string) (io.Reader, error)
}

type FileAccessor interface {
	FileWriter
	FileReader
	Exist(path string) bool
	Delete(path string) error
}

type Shell interface {
	Exec(commands []string, env *map[string]string, path []string) (string, string, int, error)
}
