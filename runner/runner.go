package runner

import (
	"fisherman/commands"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/utils"
	"os/user"
)

// Runner is main app structure
type Runner struct {
	commandList  []commands.CliCommand
	systemUser   *user.User
	config       *config.FishermanConfig
	app          *commands.AppInfo
	fileAccessor infrastructure.FileAccessor
	repository   infrastructure.Repository
}

// NewRunnerArgs is structure to pass arguments in constructor
type NewRunnerArgs struct {
	CommandList []commands.CliCommand
	Files       infrastructure.FileAccessor
	SystemUser  *user.User
	Config      *config.FishermanConfig
	ConfigInfo  *config.LoadInfo
	Cwd         string
	Executable  string
	Repository  infrastructure.Repository
}

// NewRunner is constructor for Runner
func NewRunner(args NewRunnerArgs) *Runner {
	configInfo := args.ConfigInfo

	return &Runner{
		args.CommandList,
		args.SystemUser,
		args.Config,
		&commands.AppInfo{
			Executable:         args.Executable,
			Cwd:                args.Cwd,
			GlobalConfigPath:   configInfo.GlobalConfigPath,
			LocalConfigPath:    configInfo.LocalConfigPath,
			RepoConfigPath:     configInfo.RepoConfigPath,
			IsRegisteredInPath: utils.IsCommandExists(constants.AppConfigName),
		},
		args.Files,
		args.Repository,
	}
}
