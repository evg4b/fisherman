package app

import (
	"fisherman/internal"
	"os"
	"os/signal"
)

type appOption = func(app *FishermanApp)

// WithCommands setups commands lists for application.
func WithCommands(commands []internal.CliCommand) appOption {
	return func(app *FishermanApp) {
		app.commands = commands
	}
}

// WithCwd setups current working directory (CWD) for application.
func WithCwd(cwd string) appOption {
	return func(app *FishermanApp) {
		app.cwd = cwd
	}
}

func WithSistemInterruptSignals() appOption {
	return func(app *FishermanApp) {
		app.interaption = make(chan os.Signal, 1)
		signal.Notify(app.interaption, os.Interrupt)
	}
}

func WithInterruptChanel(chanel chan os.Signal) appOption {
	return func(app *FishermanApp) {
		app.interaption = chanel
	}
}
