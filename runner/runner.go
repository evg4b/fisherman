package runner

import (
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
func NewRunner(args Args) *Runner {
	configInfo := args.ConfigInfo

	return &Runner{
		args.Commands,
		args.SystemUser,
		args.Config,
		&clicontext.AppInfo{
			Executable:       args.Executable,
			Cwd:              args.Cwd,
			GlobalConfigPath: configInfo.GlobalConfigPath,
			LocalConfigPath:  configInfo.LocalConfigPath,
			RepoConfigPath:   configInfo.RepoConfigPath,
		},
		args.Files,
		args.Repository,
		args.Shell,
	}
}
