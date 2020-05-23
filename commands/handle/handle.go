package handle

import (
	"fisherman/commands"
	"flag"
)

type Command struct {
	fs   *flag.FlagSet
	hook string
}

func NewCommand(handling flag.ErrorHandling) *Command {
	fs := flag.NewFlagSet("handle", handling)
	c := &Command{fs: fs}
	fs.StringVar(&c.hook, "hook", "", "")
	return c
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *Command) Run(ctx commands.Context) error {
	print("################")
	print(c.hook)
	print("################")
	return nil
}

func (c *Command) Name() string {
	return c.fs.Name()
}
