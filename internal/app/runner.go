package app

import (
	"context"
	c "fisherman/commands"
	i "fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/internal/appcontext"
	"fmt"
	"io"
	"strings"
)

type FishermanApp struct {
	fs       i.FileSystem
	shell    i.Shell
	repo     i.Repository
	output   io.Writer
	commands []c.CliCommand
}

func (r *FishermanApp) Run(baseCtx context.Context, args []string) error {
	if len(args) < 1 {
		log.Debug("No command detected")
		r.PrintDefaults()

		return nil
	}

	commandName := args[0]
	log.Debugf("Called command '%s'", commandName)

	for _, command := range r.commands {
		if strings.EqualFold(command.Name(), commandName) {
			err := command.Init(args[1:])
			if err != nil {
				return err
			}

			ctx := appcontext.NewContextBuilder().
				WithContext(baseCtx).
				WithFileSystem(r.fs).
				WithShell(r.shell).
				WithRepository(r.repo).
				WithArgs(args).
				WithOutput(log.InfoOutput).
				Build()

			log.Debugf("Command '%s' was initialized", commandName)
			if err := command.Run(ctx); err != nil {
				log.Debugf("Command '%s' finished with error, %v", commandName, err)

				return err
			}

			log.Debugf("Command '%s' finished witout error", commandName)

			return nil
		}
	}

	log.Debugf("Command %s not found", commandName)

	return fmt.Errorf("unknown command: %s", commandName)
}
