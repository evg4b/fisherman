package remove

import (
	"fisherman/commands"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"flag"
	"path/filepath"
)

// Command is structure for storage information about remove command
type Command struct {
	fs *flag.FlagSet
}

// NewCommand is constructor for init command
func NewCommand(handling flag.ErrorHandling) *Command {
	defer log.Debug("Remove command created")
	fs := flag.NewFlagSet("remove", handling)

	return &Command{
		fs: fs,
	}
}

// Init initialize handle command
func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

// Run executes init command
func (c *Command) Run(ctx *commands.CommandContext) error {
	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(ctx.App.Cwd, ".git", "hooks", hookName)
		if ctx.Files.Exist(hookPath) {
			log.Debugf("Hook '%s' exists", hookName)
			err := ctx.Files.Delete(hookPath)
			if err != nil {
				return err
			}

			log.Infof("Hook '%s' (%s) was removed", hookName, hookPath)
		}
	}

	return nil
}

// Name returns namand name
func (c *Command) Name() string {
	return c.fs.Name()
}
