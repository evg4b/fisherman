package infrastructure

import (
	"context"
	"fisherman/pkg/shell"
	"io"
	"os"
	"time"

	"github.com/spf13/afero"
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
	Chmod(name string, mode os.FileMode) error
	Chown(name string, uid, gid int) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
	Create(name string) (afero.File, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Name() string
	Open(name string) (afero.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (afero.File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(name string) (os.FileInfo, error)
}

type Shell interface {
	Exec(context.Context, io.Writer, string, shell.ShScript) error
}
