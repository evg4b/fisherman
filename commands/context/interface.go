package context

import (
	"fisherman/config"
	"fisherman/infrastructure/git"
	"fisherman/infrastructure/io"
	"os/user"
)

// Context is an abstract layer for accessing program data
type Context interface {
	GetGitInfo() (*git.RepositoryInfo, error)
	GetFileAccessor() io.FileAccessor
	GetCurrentUser() *user.User
	GetCwd() string
	GetAppInfo() (*AppInfo, error)
	GetConfiguration() *config.FishermanConfig
}
