package main

import (
	"context"
	"fisherman/commands"
	"fisherman/commands/handle"
	"fisherman/commands/initialize"
	"fisherman/commands/remove"
	"fisherman/commands/version"
	"fisherman/configuration"
	. "fisherman/constants" // nolint
	"fisherman/infrastructure/filesystem"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/infrastructure/vcs"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/hookfactory"
	"fisherman/internal/runner"
	"fisherman/utils"
	"io"
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
	sysShell := shell.NewShell(os.Stdout, cwd)
	repository := vcs.NewGitRepository(cwd)
	compiler := configcompiler.NewCompiler(repository, config.GlobalVariables, cwd)

	hookFactory := hookfactory.NewFactory(func(args []string, output io.Writer) *internal.Context {
		return internal.NewInternalContext(ctx, fileSystem, sysShell, repository, args, output)
	}, compiler)

	hooksHandles := hookfactory.HandlerList{
		CommitMsgHook:         hookFactory.CommitMsg(config.Hooks.CommitMsgHook),
		PreCommitHook:         hookFactory.PreCommit(config.Hooks.PreCommitHook),
		PrePushHook:           hookFactory.PrePush(config.Hooks.PrePushHook),
		PrepareCommitMsgHook:  hookFactory.PrepareCommitMsg(config.Hooks.PrepareCommitMsgHook),
		ApplyPatchMsgHook:     hookfactory.NotSupported,
		FsMonitorWatchmanHook: hookfactory.NotSupported,
		PostUpdateHook:        hookfactory.NotSupported,
		PreApplyPatchHook:     hookfactory.NotSupported,
		PreRebaseHook:         hookfactory.NotSupported,
		PreReceiveHook:        hookfactory.NotSupported,
		UpdateHook:            hookfactory.NotSupported,
	}

	appInfo := internal.AppInfo{
		Executable: executable,
		Cwd:        cwd,
		Configs:    configFiles,
	}

	commands := []commands.CliCommand{
		initialize.NewCommand(fileSystem, &appInfo, usr),
		handle.NewCommand(hooksHandles, &config.Hooks, &appInfo),
		remove.NewCommand(fileSystem, &appInfo, usr),
		version.NewCommand(),
	}

	instance := runner.NewRunner(commands, &appInfo)
	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
