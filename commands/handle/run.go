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
	err = c.header(ctx)
	if err != nil {
		return err
	}
	configuration := ctx.GetConfiguration()
	fmt.Println(configuration)
	return nil
}

func (c *Command) init(args []string) error {
	err := c.fs.Parse(args)
	if err == nil {
		c.args = c.fs.Args()
	}
	return err
}
