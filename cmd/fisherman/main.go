package main

import (
	"context"
	"fisherman/internal"
	"fisherman/internal/app"
	"fisherman/internal/commands/handle"
	"fisherman/internal/commands/initialize"
	"fisherman/internal/commands/remove"
	"fisherman/internal/commands/version"
	"fisherman/internal/configuration"
	"fisherman/internal/expression"
	"fisherman/internal/handling"
	"fisherman/internal/utils"
	"fisherman/pkg/log"
	"fisherman/pkg/shell"
	"fisherman/pkg/vcs"
	"os"
	"os/user"

	"github.com/spf13/afero"
)

const fatalExitCode = 1

func main() {
	defer utils.PanicInterceptor(os.Exit, fatalExitCode)

	usr, err := user.Current()
	utils.HandleCriticalError(err)

	cwd, err := os.Getwd()
	utils.HandleCriticalError(err)

	executable, err := os.Executable()
	utils.HandleCriticalError(err)

	fs := afero.NewOsFs()

	loader := configuration.NewLoader(usr, cwd, fs)
	configs, err := loader.FindConfigFiles()
	utils.HandleCriticalError(err)

	config, err := loader.Load(configs)
	utils.HandleCriticalError(err)

	log.Configure(config.Output)

	ctx := context.Background()
	engine := expression.NewGoExpressionEngine()
	hookFactory := handling.NewHookHandlerFactory(engine, config.Hooks)

	appInfo := internal.AppInfo{
		Executable: executable,
		Cwd:        cwd,
		Configs:    configs,
	}

	fishermanApp := app.NewAppBuilder().
		WithCommands([]internal.CliCommand{
			initialize.NewCommand(fs, appInfo, usr),
			handle.NewCommand(hookFactory, &config.Hooks, appInfo),
			remove.NewCommand(fs, appInfo, usr),
			version.NewCommand(),
		}).
		WithFs(fs).
		WithOutput(os.Stdout).
		WithRepository(vcs.NewGitRepository(cwd)).
		WithShell(shell.NewShell(os.Stdout, cwd, config.DefaultShell)).
		Build()

	if err = fishermanApp.Run(ctx, os.Args[1:]); err != nil {
		panic(err)
	}
}
