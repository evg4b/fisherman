package commands

import (
	"fisherman/config"
	"fisherman/infrastructure/git"
	"fisherman/infrastructure/io"
	"os/user"
)

type Context interface {
	GetGitInfo() (*git.RepositoryInfo, error)
	GetFileAccessor() io.FileAccessor
	GetCurrentUser() *user.User
}

type ConfigPaths struct {
	GlobalConfigPath *string
	RepoConfigPath   *string
	LocalConfigPath  *string
}

type CliCommandContext struct {
	repoInfo     *git.RepositoryInfo
	fileAccessor io.FileAccessor
	usr          *user.User
	cwd          string
	appPath      string
	config       *config.FishermanConfig
	configPaths  *ConfigPaths
}

type CliCommandContextParams struct {
	RepoInfo     *git.RepositoryInfo
	FileAccessor io.FileAccessor
	Usr          *user.User
	Cwd          string
	AppPath      string
	ConfigInfo   *config.LoadInfo
}

func NewContext(params CliCommandContextParams) *CliCommandContext {
	configInfo := params.ConfigInfo
	configPaths := ConfigPaths{
		configInfo.GlobalConfigPath,
		configInfo.RepoConfigPath,
		configInfo.LocalConfigPath,
	}
	return &CliCommandContext{
		params.RepoInfo,
		params.FileAccessor,
		params.Usr,
		params.Cwd,
		params.AppPath,
		configInfo.Config,
		&configPaths,
	}
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
