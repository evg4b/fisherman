package remove

import (
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"flag"
	"path/filepath"
)

// Command is structure for storage information about remove command
type Command struct {
	flagSet *flag.FlagSet
	usage   string
}

// NewCommand is constructor for init command
func NewCommand(handling flag.ErrorHandling) *Command {
	defer log.Debug("Remove command created")

	return &Command{
		flagSet: flag.NewFlagSet("remove", handling),
		usage:   "removes fisherman from git repository",
	}
}

// Init initialize handle command
func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

// Run executes init command
func (c *Command) Run(ctx *clicontext.CommandContext) error {
	filesToDelete := []string{
		config.BuildFileConfigPath(ctx.App.Cwd, ctx.User, config.RepoMode),
		config.BuildFileConfigPath(ctx.App.Cwd, ctx.User, config.LocalMode),
	}

	for _, hookName := range constants.HooksNames {
		filesToDelete = append(filesToDelete, filepath.Join(ctx.App.Cwd, ".git", "hooks", hookName))
	}

	for _, hookPath := range filesToDelete {
		if ctx.Files.Exist(hookPath) {
			err := ctx.Files.Delete(hookPath)
			if err != nil {
				return err
			}

			log.Infof("File '%s' was removed", hookPath)
		}
	}

	return nil
}

// Name returns command name
func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}
