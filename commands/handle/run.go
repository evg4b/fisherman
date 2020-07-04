package handle

import (
	"fisherman/commands/context"
	"fmt"
)

// Run executes handle command
func (c *Command) Run(ctx context.Context, args []string) error {
	err := c.init(args)
	if err != nil {
		return err
	}
	if hookHandler, ok := c.handlers[c.hook]; ok {
		err = c.header(ctx, c.hook)
		if err != nil {
			return err
		}
		hookHandler(ctx, ctx.GetConfiguration())
		return nil
	}
	return fmt.Errorf("%s is not valid hook name", c.hook)
}

func (c *Command) init(args []string) error {
	err := c.fs.Parse(args)
	if err == nil {
		c.args = c.fs.Args()
	}
	return err
}
