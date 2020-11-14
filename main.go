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
	"fisherman/internal/handling"
	"fisherman/internal/runner"
	"fisherman/internal/validation"
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

	extractor := validation.NewConfigExtractor(repository, conf.GlobalVariables, cwd)

	hooksHandles := map[string]handling.Handler{
		constants.CommitMsgHook:         hooks.CommitMsg(factory, conf.Hooks.CommitMsgHook, extractor),
		constants.PreCommitHook:         hooks.PreCommit(factory, conf.Hooks.PreCommitHook, extractor, sh),
		constants.ApplyPatchMsgHook:     new(handling.NotSupportedHandler),
		constants.FsMonitorWatchmanHook: new(handling.NotSupportedHandler),
		constants.PostUpdateHook:        new(handling.NotSupportedHandler),
		constants.PreApplyPatchHook:     new(handling.NotSupportedHandler),
		constants.PrePushHook:           new(handling.NotSupportedHandler),
		constants.PreRebaseHook:         new(handling.NotSupportedHandler),
		constants.PreReceiveHook:        new(handling.NotSupportedHandler),
		constants.PrepareCommitMsgHook:  new(handling.NotSupportedHandler),
		constants.UpdateHook:            new(handling.NotSupportedHandler),
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
