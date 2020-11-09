package runner

import (
	"context"
	"fisherman/clicontext"
	"fisherman/commands"
	"fisherman/config"
	"fisherman/infrastructure"
	"os/user"
)

// Runner is main app structure
type Runner struct {
	commands   []commands.CliCommand
	systemUser *user.User
	config     *config.FishermanConfig
	app        *clicontext.AppInfo
	fileSystem infrastructure.FileSystem
	repository infrastructure.Repository
	shell      infrastructure.Shell
	context    context.Context
}

// Args is structure to pass arguments in constructor
type Args struct {
	Commands   []commands.CliCommand
	Files      infrastructure.FileSystem
	Shell      infrastructure.Shell
	SystemUser *user.User
	Config     *config.FishermanConfig
	ConfigInfo *config.LoadInfo
	Cwd        string
	Executable string
	Repository infrastructure.Repository
}

// NewRunner is constructor for Runner
func NewRunner(ctx context.Context, args Args) *Runner {
	configInfo := args.ConfigInfo

	return &Runner{
		commands:   args.Commands,
		systemUser: args.SystemUser,
		config:     args.Config,
		app: &clicontext.AppInfo{
			Executable:       args.Executable,
			Cwd:              args.Cwd,
			GlobalConfigPath: configInfo.GlobalConfigPath,
			LocalConfigPath:  configInfo.LocalConfigPath,
			RepoConfigPath:   configInfo.RepoConfigPath,
		},
		fileSystem: args.Files,
		repository: args.Repository,
		shell:      args.Shell,
		context:    ctx,
	}
}
