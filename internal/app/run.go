package app

import (
	"context"
	i "fisherman/internal"
	"fisherman/internal/appcontext"
	"fisherman/pkg/log"
	"io"
)

type FishermanApp struct {
	cwd      string
	fs       i.FileSystem
	shell    i.Shell
	repo     i.Repository
	output   io.Writer
	commands CliCommands
}

func (r *FishermanApp) Run(baseCtx context.Context, args []string) error {
	if len(args) < 1 {
		log.Debug("No command detected")
		r.PrintDefaults()

		return nil
	}

	commandName := args[0]
	command, err := r.commands.GetCommand(args)
	if err != nil {
		return err
	}

	ctx := appcontext.NewContextBuilder().
		WithCwd(r.cwd).
		WithContext(baseCtx).
		WithFileSystem(r.fs).
		WithShell(r.shell).
		WithRepository(r.repo).
		WithArgs(args).
		WithOutput(log.InfoOutput).
		Build()

	if err := command.Run(ctx); err != nil {
		log.Debugf("Command '%s' finished with error, %v", commandName, err)

		return err
	}

	log.Debugf("Command '%s' finished witout error", commandName)

	return nil
}
