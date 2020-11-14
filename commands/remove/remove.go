package remove

import (
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/internal/clicontext"
	"flag"
	"path/filepath"
)

type Command struct {
	flagSet *flag.FlagSet
	usage   string
}

func NewCommand(handling flag.ErrorHandling) *Command {
	defer log.Debug("Remove command created")

	return &Command{
		flagSet: flag.NewFlagSet("remove", handling),
		usage:   "removes fisherman from git repository",
	}
}

func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

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

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}
