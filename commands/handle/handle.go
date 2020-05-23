package handle

import (
	"fisherman/commands"
	"flag"
	"os"
)

type Command struct {
	fs   *flag.FlagSet
	hook string
	args []string
}

func NewCommand(handling flag.ErrorHandling) *Command {
	fs := flag.NewFlagSet("handle", handling)
	c := &Command{fs: fs}
	fs.StringVar(&c.hook, "hook", "", "")
	return c
}

func (c *Command) Init(args []string) error {
	err := c.fs.Parse(args)
	if err == nil {
		c.args = c.fs.Args()
	}
	return err
}

func (c *Command) Run(ctx commands.Context) error {
	c.header()
	return nil
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) header() {
	info := HookInfo{
		Hook:             c.hook,
		GlobalConfigPath: "demo/",
		Version:          "0.0.1",
	}
	printHookHeader(&info, os.Stdout)
}
