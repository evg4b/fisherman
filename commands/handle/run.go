package handle

import (
	"fisherman/commands/context"
	"fmt"
	"strings"
)

// Run executes handle command
func (c *Command) Run(ctx *context.CommandContext, args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	if hookHandler, ok := c.handlers[strings.ToLower(c.hook)]; ok {
		header(ctx, c.hook)
		return hookHandler(ctx, c.fs.Args())
	}

	return fmt.Errorf("'%s' is not valid hook name", c.hook)
}
