package commands

import (
	"fisherman/config"
	"fisherman/config/hooks"
	"fisherman/infrastructure"
	"os/user"

	"github.com/imdario/mergo"
)

// CommandContext is cli context structure
type CommandContext struct {
	Config     *config.HooksConfig
	User       *user.User
	App        *AppInfo
	Files      infrastructure.FileAccessor
	Repository infrastructure.Repository
	Variables  map[string]interface{}
}

// AppInfo is application info structure
type AppInfo struct {
	Cwd                string
	Executable         string
	GlobalConfigPath   string
	LocalConfigPath    string
	RepoConfigPath     string
	IsRegisteredInPath bool
}

// CliCommandContextParams is structure for params in cli command context constructor
type CliCommandContextParams struct {
	FileAccessor infrastructure.FileAccessor
	Usr          *user.User
	App          *AppInfo
	Config       *config.FishermanConfig
	Repository   infrastructure.Repository
	Variables    map[string]interface{}
}

// NewContext constructor for cli command context
func NewContext(args CliCommandContextParams) *CommandContext {
	return &CommandContext{
		&args.Config.Hooks,
		args.Usr,
		args.App,
		args.FileAccessor,
		args.Repository,
		args.Variables,
	}
}

func (ctx *CommandContext) LoadAdditionalVariables(variables *hooks.Variables) error {
	branch, err := ctx.Repository.GetCurrentBranch()
	if err != nil {
		return err
	}

	additional, err := variables.GetFromBranch(branch)
	if err != nil {
		return err
	}

	err = mergo.MergeWithOverwrite(&ctx.Variables, additional)
	if err != nil {
		return err
	}

	return nil
}
