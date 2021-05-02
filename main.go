package main

import (
	"context"
	"fisherman/commands"
	"fisherman/commands/handle"
	"fisherman/commands/initialize"
	"fisherman/commands/remove"
	"fisherman/commands/version"
	"fisherman/configuration"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/infrastructure/vcs"
	"fisherman/internal"
	"fisherman/internal/expression"
	"fisherman/internal/handling"
	"fisherman/internal/runner"
	"fisherman/utils"
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

	fileSystem := afero.NewOsFs()

	loaded := configuration.NewLoader(usr, cwd, fileSystem)
	configFiles, err := loaded.FindConfigFiles()
	utils.HandleCriticalError(err)

	config, err := loaded.Load(configFiles)
	utils.HandleCriticalError(err)

	log.Configure(config.Output)

	ctx := context.Background()
	sysShell := shell.NewShell(os.Stdout, cwd, config.DefaultShell)
	repository := vcs.NewGitRepository(cwd)

	engine := expression.NewExpressionEngine()

	ctxFactory := internal.NewCtxFactory(ctx, fileSystem, sysShell, repository)
	hookFactory := handling.NewFactory(engine, config.Hooks)

	appInfo := internal.AppInfo{
		Executable: executable,
		Cwd:        cwd,
		Configs:    configFiles,
	}
	instance := runner.NewRunner([]commands.CliCommand{
		initialize.NewCommand(fileSystem, appInfo, usr),
		handle.NewCommand(hookFactory, ctxFactory, &config.Hooks, appInfo),
		remove.NewCommand(fileSystem, appInfo, usr),
		version.NewCommand(),
	})
	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
