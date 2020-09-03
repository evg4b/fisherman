package context

import (
	"fisherman/config"
	"fisherman/infrastructure/git"
	"os/user"
)

// GetCurrentUser returns information about currect user
func (ctx *CommandContext) GetCurrentUser() *user.User {
	return ctx.usr
}

// GetHookConfiguration returns hooks configurations
func (ctx *CommandContext) GetHookConfiguration() *config.HooksConfig {
	return &ctx.config.Hooks
}

// GetCwd returns currect working directory
func (ctx *CommandContext) GetCwd() string {
	return ctx.cwd
}

// GetGitInfo returns information about git repository
func (ctx *CommandContext) GetGitInfo() (*git.RepositoryInfo, error) {
	return git.GetRepositoryInfo("./")
}
