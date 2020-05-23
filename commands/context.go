package commands

import (
	"fisherman/infrastructure/git"
	"fisherman/infrastructure/io"
	"os/user"
)

type Context interface {
	GetGitInfo() (*git.RepositoryInfo, error)
	GetFileAccessor() io.FileAccessor
	GetCurrentUser() *user.User
}

type CliCommandContext struct {
	repoInfo     *git.RepositoryInfo
	fileAccessor io.FileAccessor
	usr          *user.User
}

type CliCommandContextParams struct {
	RepoInfo     *git.RepositoryInfo
	FileAccessor io.FileAccessor
	Usr          *user.User
	cwd          string
}

func NewContext(repoInfo *git.RepositoryInfo, fileAccessor io.FileAccessor, usr *user.User) *CliCommandContext {
	return &CliCommandContext{repoInfo, fileAccessor, usr}
}

func (ctx *CliCommandContext) GetGitInfo() (*git.RepositoryInfo, error) {
	return git.GetRepositoryInfo("./")
}

func (ctx *CliCommandContext) GetFileAccessor() io.FileAccessor {
	return ctx.fileAccessor
}

func (ctx *CliCommandContext) GetCurrentUser() *user.User {
	return ctx.usr
}
