package internal

import (
	"context"
	"fisherman/pkg/vcs"
	"io"

	"github.com/go-git/go-billy/v5"
)

// ExecutionContext is interface to access execution context.
type ExecutionContext interface {
	context.Context
	GlobalVariables() (map[string]interface{}, error)
	Files() billy.Filesystem
	Repository() Repository
	Args() []string
	Arg(index int) (string, error)
	Output() io.WriteCloser
	Message() (string, error)
	Cancel()
	Cwd() string
	Env() []string
}

// TODO: Remove this structure after option pattern implementation in internal/commands/... .
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

// CliCommand is interface to define cli command.
type CliCommand interface {
	Init(args []string) error
	Run(ctx ExecutionContext) error
	Name() string
	Description() string
}
