package internal

import (
	"context"
	"fisherman/pkg/vcs"
)

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
	Run(ctx context.Context, args []string) error
	Name() string
	Description() string
}
