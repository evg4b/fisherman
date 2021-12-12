package app

import (
	"context"
	"fisherman/internal"
	"fisherman/pkg/guards"
	"fisherman/pkg/log"
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
)

// FishermanApp is main application structure.
type FishermanApp struct {
	cwd         string
	fs          billy.Filesystem
	repo        internal.Repository
	output      io.Writer
	commands    CliCommands
	env         []string
	interaption chan os.Signal
}

// NewFishermanApp is an fisherman application constructor.
func NewFishermanApp(options ...appOption) *FishermanApp {
	app := FishermanApp{
		output:      io.Discard,
		commands:    CliCommands{},
		cwd:         "",
		env:         []string{},
		interaption: make(chan os.Signal),
	}

	for _, option := range options {
		option(&app)
	}

	guards.ShouldBeDefined(app.fs, "FileSystem should be configured")
	guards.ShouldBeDefined(app.repo, "Repository should be configured")
	guards.ShouldBeDefined(app.commands, "Commands should be configured")

	return &app
}

// Run runs fisherman application.
func (r *FishermanApp) Run(baseCtx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(baseCtx)
	subscribeInteruption(r.interaption, func() {
		log.Debug("application received interapt event")
		cancel()
	})

	if len(args) < 1 {
		log.Debug("No command detected")
		r.PrintDefaults()

		return nil
	}

	command, err := r.commands.GetCommand(args)
	if err != nil {
		return err
	}

	if err := command.Run(ctx); err != nil {
		log.Debugf("Command '%s' finished with error, %v", command.Name(), err)

		return err
	}

	log.Debugf("Command '%s' finished witout error", command.Name())

	return nil
}

func subscribeInteruption(interaption chan os.Signal, action func()) {
	go func() {
		<-interaption
		action()
	}()
}
