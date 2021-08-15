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
	"fisherman/pkg/guards"
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
	guards.NoError(err)

	cwd, err := os.Getwd()
	guards.NoError(err)

	executablePath, err := os.Executable()
	guards.NoError(err)

	fs := afero.NewOsFs()

	configLoader := configuration.NewLoader(usr, cwd, fs)

	configs, err := configLoader.FindConfigFiles()
	guards.NoError(err)

	config, err := configLoader.Load(configs)
	guards.NoError(err)

	log.Configure(config.Output)

	ctx := context.Background()
	engine := expression.NewGoExpressionEngine()

	hookFactory := handling.NewHookHandlerFactory(engine, config.Hooks)

	appInfo := internal.AppInfo{
		Executable: executablePath,
		Cwd:        cwd,
		Configs:    configs,
	}

	shell := shell.NewShell().
		WithWorkingDirectory(cwd).
		WithDefaultShell(utils.GetOrDefault(config.DefaultShell, shell.PlatformDefaultShell))

	fishermanApp := app.NewAppBuilder().
		WithCommands([]internal.CliCommand{
			initialize.NewCommand(fs, appInfo, usr),
			handle.NewCommand(hookFactory, &config.Hooks, appInfo),
			remove.NewCommand(fs, appInfo, usr),
			version.NewCommand(),
		}).
		WithFs(fs).
		WithOutput(os.Stdout).
		WithRepository(vcs.OpenGitRepository(cwd)).
		WithShell(shell).
		Build()

	if err = fishermanApp.Run(ctx, os.Args[1:]); err != nil {
		panic(err)
	}
}
