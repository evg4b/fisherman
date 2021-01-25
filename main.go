package main

import (
	"context"
	"fisherman/commands"
	"fisherman/commands/handle"
	"fisherman/commands/initialize"
	"fisherman/commands/remove"
	"fisherman/commands/version"
	"fisherman/configuration"
	"fisherman/infrastructure/filesystem"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/infrastructure/vcs"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/hookfactory"
	"fisherman/internal/runner"
	"fisherman/utils"
	"os"
	"os/user"
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

	fileSystem := filesystem.NewLocalFileSystem()

	config, configFiles, err := configuration.Load(cwd, usr, fileSystem)
	utils.HandleCriticalError(err)

	log.Configure(config.Output)

	ctx := context.Background()
	sysShell := shell.NewShell(os.Stdout, cwd, config.DefaultShell)
	repository := vcs.NewGitRepository(cwd)

	ctxFactory := internal.NewCtxFactory(ctx, fileSystem, sysShell, repository)
	extractor := configcompiler.NewConfigExtractor(repository, config.GlobalVariables, cwd)
	hookFactory := hookfactory.NewFactory(extractor, config.Hooks)

	appInfo := internal.AppInfo{
		Executable: executable,
		Cwd:        cwd,
		Configs:    configFiles,
	}

	commands := []commands.CliCommand{
		initialize.NewCommand(fileSystem, &appInfo, usr),
		handle.NewCommand(hookFactory, ctxFactory, &config.Hooks, &appInfo),
		remove.NewCommand(fileSystem, &appInfo, usr),
		version.NewCommand(),
	}

	instance := runner.NewRunner(commands, &appInfo)
	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
