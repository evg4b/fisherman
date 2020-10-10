package commands

import (
	"fisherman/config"
	"fisherman/config/hooks"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"os/user"

	"github.com/mkideal/pkg/errors"

	"github.com/imdario/mergo"
)

// CommandContext is cli context structure
type CommandContext struct {
	Config     *config.HooksConfig
	User       *user.User
	App        *AppInfo
	Files      infrastructure.FileAccessor
	Repository infrastructure.Repository
	Shell      infrastructure.Shell
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
	Shell        infrastructure.Shell
}

// NewContext constructor for cli command context
func NewContext(args CliCommandContextParams) *CommandContext {
	return &CommandContext{
		&args.Config.Hooks,
		args.Usr,
		args.App,
		args.FileAccessor,
		args.Repository,
		args.Shell,
		args.Variables,
	}
}

func (ctx *CommandContext) LoadAdditionalVariables(variables *hooks.Variables) error {
	err := ctx.load(ctx.Repository.GetCurrentBranch, variables.GetFromBranch)
	if err != nil {
		return err
	}

	err = ctx.load(ctx.Repository.GetLastTag, variables.GetFromTag)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *CommandContext) load(
	source func() (string, error),
	load func(string) (map[string]interface{}, error),
) error {
	sourceString, err := source()
	if err != nil {
		log.Debugf("Failed getting current branch: %s\n%s", err, errors.Wrap(err))

		return err
	}

	additionalValues, err := load(sourceString)
	if err != nil {
		log.Debugf("Failed getting variables from branch: %s\n%s", err, errors.Wrap(err))

		return err
	}

	err = mergo.MergeWithOverwrite(&ctx.Variables, additionalValues)
	if err != nil {
		log.Debugf("Failed merging variables: %s\n%s", err, errors.Wrap(err))

		return err
	}

	return nil
}
