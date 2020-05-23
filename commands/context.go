package commands

import (
	"fisherman/infrastructure/git"
	"fisherman/infrastructure/io"
)

type Context interface {
	GetGitInfo() (*git.RepositoryInfo, error)
	GetFileAccessor() io.FileAccessor
}

type CliCommandContext struct {
	repoInfo     *git.RepositoryInfo
	fileAccessor io.FileAccessor
}

func NewContext(repoInfo *git.RepositoryInfo, fileAccessor io.FileAccessor) *CliCommandContext {
	return &CliCommandContext{repoInfo, fileAccessor}
}

func (ctx *CliCommandContext) GetGitInfo() (*git.RepositoryInfo, error) {
	return git.GetRepositoryInfo("./")
}

func (ctx *CliCommandContext) GetFileAccessor() io.FileAccessor {
	return ctx.fileAccessor
}
