package app

import (
	"github.com/evg4b/fisherman/internal"
	"os"
	"os/signal"
)

type AppOption = func(app *FishermanApp)

// WithCommands setups commands lists for application.
func WithCommands(commands []internal.CliCommand) AppOption {
	return func(app *FishermanApp) {
		app.commands = commands
	}
}

// WithCwd setups current working directory (CWD) for application.
func WithCwd(cwd string) AppOption {
	return func(app *FishermanApp) {
		app.cwd = cwd
	}
}

func WithSistermInterruptSignals() AppOption {
	return func(app *FishermanApp) {
		app.interruption = make(chan os.Signal, 1)
		signal.Notify(app.interruption, os.Interrupt)
	}
}

func WithInterruptChanel(chanel chan os.Signal) AppOption {
	return func(app *FishermanApp) {
		app.interruption = chanel
	}
}
