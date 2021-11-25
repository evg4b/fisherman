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

	"github.com/go-git/go-billy/v5/osfs"
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

	fs := osfs.New("")

	configLoader := configuration.NewLoader(usr, cwd, fs)

	configs, err := configLoader.FindConfigFiles()
	guards.NoError(err)

	config, err := configLoader.Load(configs)
	guards.NoError(err)

	log.Configure(config.Output)

	// TODO: Add Signal Interrupt handler
	ctx := context.Background()
	engine := expression.NewGoExpressionEngine()

	hookFactory := handling.NewHookHandlerFactory(engine, config.Hooks)

	appInfo := internal.AppInfo{
		Executable: executablePath,
		Cwd:        cwd,
		Configs:    configs,
	}

	defaultShell := utils.FirstNotEmpty(config.DefaultShell, shell.PlatformDefaultShell)
	shell := shell.NewShell(
		shell.WithWorkingDirectoryOld(cwd),
		shell.WithDefaultShell(defaultShell),
		shell.WithEnvOld(os.Environ()),
	)

	fishermanApp := app.NewFishermanApp(
		app.WithCommands([]internal.CliCommand{
			initialize.NewCommand(fs, appInfo, usr),
			handle.NewCommand(hookFactory, &config.Hooks, appInfo),
			remove.NewCommand(fs, appInfo, usr),
			version.NewCommand(),
		}),
		app.WithFs(fs),
		app.WithOutput(os.Stdout),
		app.WithRepository(vcs.NewRepository(
			vcs.WithFsRepo(cwd),
		)),
		app.WithShell(shell),
		app.WithEnv(os.Environ()),
	)

	// TODO: Add interrupt event hadling (stopping)
	if err = fishermanApp.Run(ctx, os.Args[1:]); err != nil {
		panic(err)
	}
}
