package runner

import (
	"fisherman/clicontext"
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
	app          *clicontext.AppInfo
	fileAccessor infrastructure.FileAccessor
	repository   infrastructure.Repository
	shell        infrastructure.Shell
}

// Args is structure to pass arguments in constructor
type Args struct {
	CommandList []commands.CliCommand
	Files       infrastructure.FileAccessor
	Shell       infrastructure.Shell
	SystemUser  *user.User
	Config      *config.FishermanConfig
	ConfigInfo  *config.LoadInfo
	Cwd         string
	Executable  string
	Repository  infrastructure.Repository
}

// NewRunner is constructor for Runner
func NewRunner(args Args) *Runner {
	configInfo := args.ConfigInfo

	return &Runner{
		args.CommandList,
		args.SystemUser,
		args.Config,
		&clicontext.AppInfo{
			Executable:         args.Executable,
			Cwd:                args.Cwd,
			GlobalConfigPath:   configInfo.GlobalConfigPath,
			LocalConfigPath:    configInfo.LocalConfigPath,
			RepoConfigPath:     configInfo.RepoConfigPath,
			IsRegisteredInPath: utils.IsCommandExists(constants.AppName),
		},
		args.Files,
		args.Repository,
		args.Shell,
	}
}
