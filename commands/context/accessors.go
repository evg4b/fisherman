package context

import (
	"fisherman/infrastructure/git"
)

// GetGitInfo returns information about git repository
func (ctx *CommandContext) GetGitInfo() (*git.RepositoryInfo, error) {
	return git.GetRepositoryInfo("./")
}
