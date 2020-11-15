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

	appPath, err := os.Executable()
	utils.HandleCriticalError(err)

	conf, configInfo, err := config.Load(cwd, usr, filesystem.NewLocalFileSystem())
	utils.HandleCriticalError(err)

	log.Configure(conf.Output)

	ctx := context.Background()

	files := filesystem.NewLocalFileSystem()
	sh := shell.NewShell(os.Stdout)
	repository := vcs.NewGitRepository(cwd)

	factory := func(args []string, output io.Writer) *internal.Context {
		return internal.NewInternalContext(ctx, files, sh, repository, args, output)
	}

	compiler := configcompiler.NewCompiler(repository, conf.GlobalVariables, cwd)

	hooksHandles := hooks.HandlerList{
		constants.CommitMsgHook:         hooks.CommitMsg(factory, conf.Hooks.CommitMsgHook, compiler),
		constants.PreCommitHook:         hooks.PreCommit(factory, conf.Hooks.PreCommitHook, sh, compiler),
		constants.ApplyPatchMsgHook:     hooks.NotSupported,
		constants.FsMonitorWatchmanHook: hooks.NotSupported,
		constants.PostUpdateHook:        hooks.NotSupported,
		constants.PreApplyPatchHook:     hooks.NotSupported,
		constants.PrePushHook:           hooks.PrePush(factory, conf.Hooks.PrePushHook, sh, compiler),
		constants.PreRebaseHook:         hooks.NotSupported,
		constants.PreReceiveHook:        hooks.NotSupported,
		constants.PrepareCommitMsgHook:  hooks.PrepareCommitMsg(factory, conf.Hooks.PrepareCommitMsgHook, compiler),
		constants.UpdateHook:            hooks.NotSupported,
	}

	appInfo := internal.AppInfo{
		Executable:       appPath,
		Cwd:              cwd,
		GlobalConfigPath: configInfo.GlobalConfigPath,
		LocalConfigPath:  configInfo.LocalConfigPath,
		RepoConfigPath:   configInfo.RepoConfigPath,
	}

	instance := runner.NewRunner(
		[]commands.CliCommand{
			initialize.NewCommand(files, &appInfo, usr),
			handle.NewCommand(hooksHandles, &conf.Hooks, &appInfo),
			remove.NewCommand(files, &appInfo, usr),
			version.NewCommand(),
		},
		conf,
		&appInfo)

	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
