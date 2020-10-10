package main

import (
	"fisherman/commands"
	"fisherman/commands/handle"
	initc "fisherman/commands/init"
	"fisherman/commands/remove"
	"fisherman/commands/version"
	"fisherman/config"
	"fisherman/infrastructure/fs"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/vcs"
	"fisherman/runner"
	"fisherman/utils"
	"flag"
	"os"
	"os/user"
)

const fatalExitCode = 1
const applicationErrorCode = 2
const successCode = 0

func main() {
	defer panicInterceptor()

	usr, err := user.Current()
	utils.HandleCriticalError(err)

	cwd, err := os.Getwd()
	utils.HandleCriticalError(err)

	appPath, err := os.Executable()
	utils.HandleCriticalError(err)

	fileAccessor := fs.NewAccessor()

	conf, configInfo := config.LoadConfig(cwd, usr, fileAccessor)

	repo, err := vcs.NewGitRepository(cwd)
	utils.HandleCriticalError(err)

	log.Configure(conf.Output)
	runnerInstance := runner.NewRunner(runner.NewRunnerArgs{
		CommandList: []commands.CliCommand{
			initc.NewCommand(flag.ExitOnError),
			handle.NewCommand(flag.ExitOnError),
			remove.NewCommand(flag.ExitOnError),
			version.NewCommand(flag.ExitOnError),
		},
		Config:     conf,
		ConfigInfo: configInfo,
		Files:      fileAccessor,
		SystemUser: usr,
		Cwd:        cwd,
		Executable: appPath,
		Repository: repo,
	})

	if err = runnerInstance.Run(os.Args[1:]); err != nil {
		log.Error(err)
		exit(applicationErrorCode)
	}

	exit(successCode)
}

func panicInterceptor() {
	if err := recover(); err != nil {
		log.Errorf("Fatal error: %s", err)
		exit(fatalExitCode)
	}
}

func exit(code int) {
	log.Debugf("Process exit with code %d", code)
	os.Exit(code)
}
