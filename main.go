package main

import (
	"fisherman/commands"
	"fisherman/commands/handle"
	initc "fisherman/commands/init"
	"fisherman/config"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"fisherman/runner"
	"fisherman/utils"
	"flag"
	"os"
	"os/user"
)

const fatalExitCode = 1
const applicationErrorCode = 2

func main() {
	defer panicInterceptor()

	usr, err := user.Current()
	utils.HandleCriticalError(err)

	cwd, err := os.Getwd()
	utils.HandleCriticalError(err)

	appPath, err := os.Executable()
	utils.HandleCriticalError(err)

	fileAccessor := io.NewFileAccessor()

	conf, configInfo := config.LoadConfig(cwd, usr, fileAccessor)

	logger.Configure(conf.Output)
	runnerInstance := runner.NewRunner(runner.NewRunnerArgs{
		CommandList: []commands.CliCommand{
			initc.NewCommand(flag.ExitOnError),
			handle.NewCommand(flag.ExitOnError, fileAccessor),
		},
		Config:     conf,
		ConfigInfo: configInfo,
		Files:      fileAccessor,
		SystemUser: usr,
		Cwd:        cwd,
		Executable: appPath,
	})

	if err = runnerInstance.Run(os.Args[1:]); err != nil {
		logger.Error(err)
		os.Exit(applicationErrorCode)
	}
}

func panicInterceptor() {
	if err := recover(); err != nil {
		logger.Error("Fatal error:")
		logger.Error(err)
		os.Exit(fatalExitCode)
	}
}
