package commands

import (
	"flag"
	"fmt"
)

type InitCommand struct {
	fs *flag.FlagSet
	mode string
}

func NewInitCommand() *InitCommand {
	gc := &InitCommand{
		fs: flag.NewFlagSet("init", flag.ExitOnError),
	}
	gc.fs.StringVar(&gc.mode, "mode", "repo", "(local,repo,global)")
	return gc
}

func (g *InitCommand) Name() string {
	return g.fs.Name()
}

func (g *InitCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *InitCommand) Run() error {
	fmt.Println("Hello", g.mode, "!")
	return nil
}
