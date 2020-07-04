package handle

import (
	"fisherman/commands"
	"flag"
	"fmt"
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
	c.header(ctx)

	dd := ctx.GetConfiguration()
	fmt.Println(dd)

	return nil
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) header(ctx commands.Context) {
	app := ctx.GetAppInfo()
	info := HookInfo{
		Hook:             c.hook,
		GlobalConfigPath: formatNA(app.GlobalConfigPath),
		LocalConfigPath:  formatNA(app.LocalConfigPath),
		RepoConfigPath:   formatNA(app.RepoConfigPath),
		Version:          "0.0.1",
	}
	printHookHeader(&info, os.Stdout)
}

func formatNA(path *string) string {
	if path == nil {
		return "N/A"
	}
	return *path
}
