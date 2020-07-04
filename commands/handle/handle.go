package handle

import (
	"fisherman/infrastructure/reporter"
	"flag"
)

// Command is structure for storage information about handle command
type Command struct {
	fs       *flag.FlagSet
	hook     string
	args     []string
	reporter reporter.Reporter
}

// NewCommand is constructor for handle command
func NewCommand(handling flag.ErrorHandling, reporter reporter.Reporter) *Command {
	fs := flag.NewFlagSet("handle", handling)
	c := &Command{fs: fs, reporter: reporter}
	fs.StringVar(&c.hook, "hook", "", "")
	return c
}

// Name returns command name
func (c *Command) Name() string {
	return c.fs.Name()
}
