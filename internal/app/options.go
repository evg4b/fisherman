package app

import (
	i "fisherman/internal"
	"io"

	"github.com/go-git/go-billy/v5"
)

type appOption = func(app *FishermanApp)

func WithCommands(commands []i.CliCommand) appOption {
	return func(app *FishermanApp) {
		app.commands = commands
	}
}

func WithCwd(cwd string) appOption {
	return func(app *FishermanApp) {
		app.cwd = cwd
	}
}

func WithFs(fs billy.Filesystem) appOption {
	return func(app *FishermanApp) {
		app.fs = fs
	}
}

func WithOutput(output io.Writer) appOption {
	return func(app *FishermanApp) {
		app.output = output
	}
}

func WithShell(shell i.Shell) appOption {
	return func(app *FishermanApp) {
		app.shell = shell
	}
}

func WithRepository(repo i.Repository) appOption {
	return func(app *FishermanApp) {
		app.repo = repo
	}
}
