package internal

import (
	"context"
	"fisherman/pkg/shell"
	"io"

	"github.com/go-git/go-billy/v5"
)

type ExecutionContext interface {
	context.Context
	GlobalVariables() (map[string]interface{}, error)
	Files() billy.Filesystem
	Shell() Shell
	Repository() Repository
	Args() []string
	Arg(index int) (string, error)
	Output() io.Writer
	Message() (string, error)
	Cancel()
}

type AppInfo struct {
	Cwd        string
	Executable string
	Configs    map[string]string
}

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

type Shell interface {
	Exec(context.Context, io.Writer, string, *shell.Script) error
}

type CliCommand interface {
	Init(args []string) error
	Run(ctx ExecutionContext) error
	Name() string
	Description() string
}
