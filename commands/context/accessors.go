package context

import (
	"fisherman/config"
	"fisherman/infrastructure/git"
)

// GetHookConfiguration returns hooks configurations
func (ctx *CommandContext) GetHookConfiguration() *config.HooksConfig {
	return &ctx.config.Hooks
}

// GetGitInfo returns information about git repository
func (ctx *CommandContext) GetGitInfo() (*git.RepositoryInfo, error) {
	return git.GetRepositoryInfo("./")
}
