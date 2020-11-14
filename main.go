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
	c "fisherman/constants"
	"fisherman/infrastructure/filesystem"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/infrastructure/vcs"
	"fisherman/internal"
	"fisherman/internal/handling"
	"fisherman/internal/runner"
	"fisherman/internal/validation"
	"fisherman/utils"
	"flag"
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

	fs := filesystem.NewLocalFileSystem()
	sh := shell.NewShell(os.Stdout)
	repo := vcs.NewGitRepository(cwd)

	factory := func(args []string, output io.Writer) *internal.Context {
		return internal.NewInternalContext(ctx, fs, sh, repo, args, output)
	}

	extractor := validation.NewConfigExtractor(repo, conf.GlobalVariables, cwd)

	hooksHandles := map[string]handling.Handler{
		c.CommitMsgHook:         hooks.CommitMsg(factory, conf.Hooks.CommitMsgHook, extractor),
		c.PreCommitHook:         hooks.PreCommit(factory, conf.Hooks.PreCommitHook, extractor, sh),
		c.ApplyPatchMsgHook:     new(handling.NotSupportedHandler),
		c.FsMonitorWatchmanHook: new(handling.NotSupportedHandler),
		c.PostUpdateHook:        new(handling.NotSupportedHandler),
		c.PreApplyPatchHook:     new(handling.NotSupportedHandler),
		c.PrePushHook:           new(handling.NotSupportedHandler),
		c.PreRebaseHook:         new(handling.NotSupportedHandler),
		c.PreReceiveHook:        new(handling.NotSupportedHandler),
		c.PrepareCommitMsgHook:  new(handling.NotSupportedHandler),
		c.UpdateHook:            new(handling.NotSupportedHandler),
	}

	instance := runner.NewRunner(ctx, runner.Args{
		Commands: []commands.CliCommand{
			initialize.NewCommand(flag.ExitOnError),
			handle.NewCommand(flag.ExitOnError, hooksHandles),
			remove.NewCommand(flag.ExitOnError),
			version.NewCommand(flag.ExitOnError),
		},
		Config:     conf,
		ConfigInfo: configInfo,
		Files:      fs,
		SystemUser: usr,
		Cwd:        cwd,
		Executable: appPath,
		Repository: repo,
		Shell:      sh,
	})

	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
