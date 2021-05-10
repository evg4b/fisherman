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

	fs := afero.NewOsFs()

	loader := configuration.NewLoader(usr, cwd, fs)
	configFiles, err := loader.FindConfigFiles()
	utils.HandleCriticalError(err)

	config, err := loader.Load(configFiles)
	utils.HandleCriticalError(err)

	log.Configure(config.Output)

	ctx := context.Background()
	sysShell := shell.NewShell(os.Stdout, cwd, config.DefaultShell)
	repo := vcs.NewGitRepository(cwd)

	engine := expression.NewExpressionEngine()

	ctxFactory := internal.NewCtxFactory(ctx, fs, sysShell, repo)
	hookFactory := handling.NewFactory(engine, config.Hooks)

	appInfo := internal.AppInfo{
		Executable: executable,
		Cwd:        cwd,
		Configs:    configFiles,
	}
	instance := runner.NewRunner([]commands.CliCommand{
		initialize.NewCommand(fs, appInfo, usr),
		handle.NewCommand(hookFactory, ctxFactory, &config.Hooks, appInfo),
		remove.NewCommand(fs, appInfo, usr),
		version.NewCommand(),
	})
	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
