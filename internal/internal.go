package internal

import (
	"context"
	"fisherman/pkg/shell"
	"fisherman/pkg/vcs"
	"io"

	"github.com/go-git/go-billy/v5"
)

// ExecutionContext is interface to access execution context.
type ExecutionContext interface {
	context.Context
	GlobalVariables() (map[string]interface{}, error)
	Files() billy.Filesystem
	Shell() Shell
	Repository() Repository
	Args() []string
	Arg(index int) (string, error)
	Output() io.WriteCloser
	Message() (string, error)
	Cancel()
	Cwd() string
}

type AppInfo struct {
	Cwd        string
	Executable string
	Configs    map[string]string
}

// Repository is interface to comunicate with git.
type Repository interface {
	GetCurrentBranch() (string, error)
	GetUser() (vcs.User, error)
	GetLastTag() (string, error)
	AddGlob(glob string) error
	RemoveGlob(glob string) error
	GetFilesInIndex() ([]string, error)
	GetIndexChanges() (map[string]vcs.Changes, error)
}

// Shell is interface to comunicate with system shell (cmd, powersell, bash and etc.).
type Shell interface {
	Exec(context.Context, io.Writer, string, *shell.Script) error
}

// CliCommand is interface to define cli command.
type CliCommand interface {
	Init(args []string) error
	Run(ctx ExecutionContext) error
	Name() string
	Description() string
}
