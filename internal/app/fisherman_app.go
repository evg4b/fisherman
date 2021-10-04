package app

import (
	"context"
	i "fisherman/internal"
	"fisherman/internal/appcontext"
	"fisherman/pkg/guards"
	"fisherman/pkg/log"
	"io"

	"github.com/go-git/go-billy/v5"
)

type FishermanApp struct {
	cwd      string
	fs       billy.Filesystem
	shell    i.Shell
	repo     i.Repository
	output   io.Writer
	commands CliCommands
}

func NewFishermanApp(options ...appOption) *FishermanApp {
	app := FishermanApp{
		output:   io.Discard,
		commands: make(CliCommands, 0),
		cwd:      "",
	}

	for _, option := range options {
		option(&app)
	}

	guards.ShouldBeDefined(app.fs, "FileSystem should be configured")
	guards.ShouldBeDefined(app.shell, "Shell should be configured")
	guards.ShouldBeDefined(app.repo, "Repository should be configured")
	guards.ShouldBeDefined(app.commands, "Commands should be configured")

	return &app
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

	ctx := appcontext.NewContext(
		appcontext.WithCwd(r.cwd),
		appcontext.WithBaseContext(baseCtx),
		appcontext.WithFileSystem(r.fs),
		appcontext.WithShell(r.shell),
		appcontext.WithRepository(r.repo),
		appcontext.WithArgs(args),
		appcontext.WithOutput(log.InfoOutput),
	)

	if err := command.Run(ctx); err != nil {
		log.Debugf("Command '%s' finished with error, %v", commandName, err)

		return err
	}

	log.Debugf("Command '%s' finished witout error", commandName)

	return nil
}