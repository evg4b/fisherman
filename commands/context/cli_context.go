package context

import (
	"fisherman/config"
	"fisherman/infrastructure/git"
	"fisherman/infrastructure/io"
	"os/user"
)

// CliCommandContext is cli context structure
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

// CliCommandContextParams is structure for params in cli command context constructor
type CliCommandContextParams struct {
	RepoInfo     *git.RepositoryInfo
	FileAccessor io.FileAccessor
	Usr          *user.User
	Cwd          string
	AppPath      string
	ConfigInfo   *config.LoadInfo
	Path         string
}

// NewContext constructor for cli command context
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
