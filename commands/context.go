package commands

import (
	"fisherman/config"
	"fisherman/infrastructure/git"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/path"
	"os/user"
)

type Context interface {
	GetGitInfo() (*git.RepositoryInfo, error)
	GetFileAccessor() io.FileAccessor
	GetCurrentUser() *user.User
	GetCwd() string
	GetAppInfo() *AppInfo
	GetConfiguration() *config.FishermanConfig
}

type AppInfo struct {
	AppPath            string
	IsRegisteredInPath bool
	GlobalConfigPath   *string
	RepoConfigPath     *string
	LocalConfigPath    *string
}

type CliCommandContext struct {
	repoInfo         *git.RepositoryInfo
	fileAccessor     io.FileAccessor
	usr              *user.User
	cwd              string
	config           *config.FishermanConfig
	appPath          string
	globalConfigPath *string
	repoConfigPath   *string
	localConfigPath  *string
	path             string
}

func (ctx *CliCommandContext) GetConfiguration() *config.FishermanConfig {
	return ctx.config
}

func (ctx *CliCommandContext) GetAppInfo() *AppInfo {
	isRegistered, err := path.IsRegisteredInPath(ctx.path, ctx.appPath)
	if err != nil {
		// TODO Add correct error handling
		panic(err)
	}
	return &AppInfo{
		GlobalConfigPath:   ctx.globalConfigPath,
		LocalConfigPath:    ctx.localConfigPath,
		RepoConfigPath:     ctx.repoConfigPath,
		IsRegisteredInPath: isRegistered,
		AppPath:            ctx.appPath,
	}
}

type CliCommandContextParams struct {
	RepoInfo     *git.RepositoryInfo
	FileAccessor io.FileAccessor
	Usr          *user.User
	Cwd          string
	AppPath      string
	ConfigInfo   *config.LoadInfo
	Path         string
}

func NewContext(params CliCommandContextParams) *CliCommandContext {
	configInfo := params.ConfigInfo
	return &CliCommandContext{
		params.RepoInfo,
		params.FileAccessor,
		params.Usr,
		params.Cwd,
		configInfo.Config,
		params.AppPath,
		configInfo.GlobalConfigPath,
		configInfo.RepoConfigPath,
		configInfo.LocalConfigPath,
		params.Path,
	}
}

func (ctx *CliCommandContext) GetCwd() string {
	return ctx.cwd
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
