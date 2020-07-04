package context

import (
	"fisherman/config"
	"fisherman/infrastructure/git"
	"fisherman/infrastructure/io"
	"os/user"
)

// GetCurrentUser returns information about currect user
func (ctx *CliCommandContext) GetCurrentUser() *user.User {
	return ctx.usr
}

// GetConfiguration returns fisherman configurations
func (ctx *CliCommandContext) GetConfiguration() *config.FishermanConfig {
	return ctx.config
}

// GetCwd returns currect working directory
func (ctx *CliCommandContext) GetCwd() string {
	return ctx.cwd
}

// GetGitInfo returns information about git repository
func (ctx *CliCommandContext) GetGitInfo() (*git.RepositoryInfo, error) {
	return git.GetRepositoryInfo("./")
}

// GetFileAccessor returns abstract file accessor
func (ctx *CliCommandContext) GetFileAccessor() io.FileAccessor {
	return ctx.fileAccessor
}
