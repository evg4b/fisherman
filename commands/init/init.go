package init

import (
	"fisherman/commands"
	"flag"
)

type Command struct {
	fs       *flag.FlagSet
	mode     string
	absolute bool
	force    bool
}

func NewCommand(handling flag.ErrorHandling) *Command {
	fs := flag.NewFlagSet("init", handling)
	c := &Command{fs: fs}
	fs.StringVar(&c.mode, "mode", "repo", "(local,repo,global)")
	fs.BoolVar(&c.force, "force", false, "")
	fs.BoolVar(&c.force, "absolute", false, "")
	return c
}

func (c *Command) Run(ctx commands.Context) error {
	accessor := ctx.GetFileAccessor()
	info, err := ctx.GetGitInfo()
	if err != nil {
		return err
	}

	err = WriteHooks(info.Path, accessor, c.force)
	if err != nil {
		return err
	}

	err = WriteFishermanConfig(info.Path, ctx.GetCurrentUser(), c.mode, accessor)
	if err != nil {
		return err
	}

	return nil
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}
