package infrastructure

import (
	"context"
	"fisherman/infrastructure/shell"
	"io"
	"os"
	"os/user"
)

type User struct {
	UserName string
	Email    string
}

type Repository interface {
	GetCurrentBranch() (string, error)
	GetUser() (User, error)
	GetLastTag() (string, error)
	AddGlob(glob string) error
	RemoveGlob(glob string) error
	GetFilesInIndex() ([]string, error)
}

type FileSystem interface {
	Write(path, content string) error
	Read(path string) (string, error)
	Reader(path string) (io.ReadCloser, error)
	Exist(path string) bool
	Delete(path string) error
	Chmod(path string, mode os.FileMode) error
	Chown(name string, user *user.User) error
}

type Shell interface {
	Exec(context.Context, io.Writer, string, shell.ShScript) error
}
