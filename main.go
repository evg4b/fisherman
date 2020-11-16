package main

import (
	"context"
	"fisherman/commands"
	"fisherman/commands/handle"
	"fisherman/commands/handle/hooks"
	"fisherman/commands/initialize"
	"fisherman/commands/remove"
	"fisherman/commands/version"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure/filesystem"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/infrastructure/vcs"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
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

	configuration, configInfo, err := config.Load(cwd, usr, filesystem.NewLocalFileSystem())
	utils.HandleCriticalError(err)

	log.Configure(configuration.Output)

	ctx := context.Background()
	fileSystem := filesystem.NewLocalFileSystem()
	sysShell := shell.NewShell(os.Stdout, cwd)
	repository := vcs.NewGitRepository(cwd)
	compiler := configcompiler.NewCompiler(repository, configuration.GlobalVariables, cwd)

	factory := func(args []string, output io.Writer) *internal.Context {
		return internal.NewInternalContext(ctx, fileSystem, sysShell, repository, args, output)
	}

	hooksHandles := hooks.HandlerList{
		constants.CommitMsgHook: hooks.CommitMsg(
			factory,
			configuration.Hooks.CommitMsgHook,
			compiler,
		),
		constants.PreCommitHook: hooks.PreCommit(
			factory,
			configuration.Hooks.PreCommitHook,
			sysShell,
			compiler,
		),
		constants.PrePushHook: hooks.PrePush(
			factory,
			configuration.Hooks.PrePushHook,
			sysShell,
			compiler,
		),
		constants.PrepareCommitMsgHook: hooks.PrepareCommitMsg(
			factory,
			configuration.Hooks.PrepareCommitMsgHook,
			compiler,
		),
		constants.ApplyPatchMsgHook:     hooks.NotSupported,
		constants.FsMonitorWatchmanHook: hooks.NotSupported,
		constants.PostUpdateHook:        hooks.NotSupported,
		constants.PreApplyPatchHook:     hooks.NotSupported,
		constants.PreRebaseHook:         hooks.NotSupported,
		constants.PreReceiveHook:        hooks.NotSupported,
		constants.UpdateHook:            hooks.NotSupported,
	}

	appInfo := internal.AppInfo{
		Executable:       executable,
		Cwd:              cwd,
		GlobalConfigPath: configInfo.GlobalConfigPath,
		LocalConfigPath:  configInfo.LocalConfigPath,
		RepoConfigPath:   configInfo.RepoConfigPath,
	}

	commands := []commands.CliCommand{
		initialize.NewCommand(fileSystem, &appInfo, usr),
		handle.NewCommand(hooksHandles, &configuration.Hooks, &appInfo),
		remove.NewCommand(fileSystem, &appInfo, usr),
		version.NewCommand(),
	}

	instance := runner.NewRunner(commands, &appInfo)
	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
