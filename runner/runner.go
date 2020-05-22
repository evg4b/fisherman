package runner

import (
	"fisherman/commands"
	"flag"
	"fmt"
)

func Run(args []string) error {
	if len(args) < 1 {
		flag.PrintDefaults()
	}

	cmds := []commands.CliCommand{
		commands.NewInitCommand(),
	}

	subcommand := args[0]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			if err := cmd.Init(args[1:]); err != nil {
				return err
			}
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}
