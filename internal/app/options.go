package app

import (
	"fisherman/internal"
	"io"

	"github.com/go-git/go-billy/v5"
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

// WithCwd setups file system abstraction object.
func WithFs(fs billy.Filesystem) appOption {
	return func(app *FishermanApp) {
		app.fs = fs
	}
}

// WithOutput setups default output writer.
func WithOutput(output io.Writer) appOption {
	return func(app *FishermanApp) {
		app.output = output
	}
}

// WithOutput setups shell abstraction object.
func WithShell(shell internal.Shell) appOption {
	return func(app *FishermanApp) {
		app.shell = shell
	}
}

// WithRepository setups git repository abstraction object.
func WithRepository(repo internal.Repository) appOption {
	return func(app *FishermanApp) {
		app.repo = repo
	}
}

// WithEnv setups environment variables for fisherman application.
func WithEnv(env []string) appOption {
	return func(ac *FishermanApp) {
		ac.env = env
	}
}
