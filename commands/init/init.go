package init

import (
	"fisherman/commands"
	"fisherman/config"
	"flag"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
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
	info, err := ctx.GetGitInfo()
	if err != nil {
		return err
	}

	err = WriteHooks(info.Path, ctx.GetFileAccessor(), c.force)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config.DefaultConfig)
	err = ioutil.WriteFile(".fisherman.yml", data, 0777)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}
