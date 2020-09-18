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

	"github.com/fatih/color"
)

const fatalExitCode = 1
const applicationErrorCode = 2

func main() {
	defer panicInterceptor()
	fileAccessor := io.NewFileAccessor()
	usr, err := user.Current()
	utils.HandleCriticalError(err)

	cwd, err := os.Getwd()
	utils.HandleCriticalError(err)

	conf, configInfo, err := config.LoadConfig(cwd, usr, fileAccessor)
	utils.HandleCriticalError(err)

	loggerInstance := logger.NewConsoleLogger(conf.Output)
	runnerInstance := runner.NewRunner(runner.CreateRunnerArgs{
		CommandList: []commands.CliCommand{
			initc.NewCommand(flag.ExitOnError),
			handle.NewCommand(flag.ExitOnError, fileAccessor),
		},
		Config:       conf,
		ConfigInfo:   configInfo,
		FileAccessor: fileAccessor,
		Logger:       loggerInstance,
		SystemUser:   usr,
	})

	if err = runnerInstance.Run(conf, os.Args); err != nil {
		loggerInstance.Error(err)
		os.Exit(applicationErrorCode)
	}
}

func panicInterceptor() {
	if err := recover(); err != nil {
		print := color.New(color.BgRed, color.FgWhite).PrintlnFunc()
		print("Fatal error:")
		print(err)
		os.Exit(fatalExitCode)
	}
}
