package app

import (
	"context"
	"fisherman/pkg/guards"
	"fisherman/pkg/log"
	"os"
)

// FishermanApp is main application structure.
type FishermanApp struct {
	cwd          string
	commands     CliCommands
	interruption chan os.Signal
}

// NewFishermanApp is an fisherman application constructor.
func NewFishermanApp(options ...appOption) *FishermanApp {
	app := FishermanApp{
		commands:     CliCommands{},
		cwd:          "",
		interruption: make(chan os.Signal),
	}

	for _, option := range options {
		option(&app)
	}

	guards.ShouldBeDefined(app.commands, "Commands should be configured")

	return &app
}

// Run runs fisherman application.
func (r *FishermanApp) Run(baseCtx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(baseCtx)
	subscribeInteruption(r.interruption, func() {
		log.Debug("application received interapt event")
		cancel()
	})

	if len(args) < 1 {
		log.Debug("No command detected")
		r.PrintDefaults()

		return nil
	}

	commandName, commandArgs := splitArgs(args)
	command, err := r.commands.GetCommand(commandName)
	if err != nil {
		return err
	}

	if err := command.Run(ctx, commandArgs); err != nil {
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

func splitArgs(args []string) (string, []string) {
	return args[0], args[1:]
}
