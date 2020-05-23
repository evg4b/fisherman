package runner

import (
	"fisherman/commands"
	"fisherman/commands/handle"
	initc "fisherman/commands/init"
	"fisherman/infrastructure/git"
	"fisherman/infrastructure/io"
	"flag"
	"fmt"
	"os"
	"os/user"
)

type Runner struct {
	fileAccessor io.FileAccessor
	usr          *user.User
}

func NewRunner(fileAccessor io.FileAccessor, usr *user.User) *Runner {
	return &Runner{fileAccessor, usr}
}

func (runner *Runner) Run(args []string) error {
	if len(args) < 1 {
		flag.PrintDefaults()
	}

	errorHandlingMode := flag.ExitOnError
	cmds := []commands.CliCommand{
		initc.NewCommand(errorHandlingMode),
		handle.NewCommand(errorHandlingMode),
	}

	subcommand := args[0]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			if err := cmd.Init(args[1:]); err != nil {
				return err
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			info, err := git.GetRepositoryInfo(cwd)
			if err != nil {
				return err
			}

			context := commands.NewContext(info, runner.fileAccessor, runner.usr)
			return cmd.Run(context)
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}
