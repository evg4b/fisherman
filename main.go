package main

import (
	"fisherman/commands"
	"fisherman/commands/handle"
	"fisherman/commands/initialize"
	"fisherman/commands/remove"
	"fisherman/commands/version"
	"fisherman/config"
	"fisherman/infrastructure/filesystem"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/infrastructure/vcs"
	"fisherman/runner"
	"fisherman/utils"
	"flag"
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

	handling := flag.ExitOnError

	instance := runner.NewRunner(runner.Args{
		Commands: []commands.CliCommand{
			initialize.NewCommand(handling),
			handle.NewCommand(handling),
			remove.NewCommand(handling),
			version.NewCommand(handling),
		},
		Config:     conf,
		ConfigInfo: configInfo,
		Files:      filesystem.NewLocalFileSystem(),
		SystemUser: usr,
		Cwd:        cwd,
		Executable: appPath,
		Repository: vcs.NewGitRepository(cwd),
		Shell:      shell.NewShell(),
	})

	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
