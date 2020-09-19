package init

import (
	"fisherman/commands/context"
	"fisherman/config"
	"flag"
	"fmt"
)

// Command is structure for storage information about init command
type Command struct {
	fs       *flag.FlagSet
	mode     string
	absolute bool
	force    bool
}

// NewCommand is constructor for init command
func NewCommand(handling flag.ErrorHandling) *Command {
	fs := flag.NewFlagSet("init", handling)
	c := &Command{fs: fs}
	modeMessage := fmt.Sprintf("(%s, %s, %s)", config.LocalMode, config.RepoMode, config.GlobalMode)
	fs.StringVar(&c.mode, "mode", config.RepoMode, modeMessage)
	fs.BoolVar(&c.force, "force", false, "")
	fs.BoolVar(&c.force, "absolute", false, "")
	return c
}

// Run executes init command
func (c *Command) Run(ctx *context.CommandContext, args []string) error {
	c.fs.Parse(args)

	err := writeHooks(&ctx.AppInfo, ctx.FileAccessor, c.force)
	if err != nil {
		return err
	}

	configPath, err := config.BuildFileConfigPath(ctx.AppInfo.Cwd, ctx.User, c.mode)
	if err != nil {
		return err
	}

	err = writeFishermanConfig(ctx.FileAccessor, configPath)
	if err != nil {
		return err
	}

	return nil
}

// Name returns namand name
func (c *Command) Name() string {
	return c.fs.Name()
}
