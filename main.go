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
	"log"
	"os"
	"os/user"

	"github.com/fatih/color"
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

	conf, configInfo, err := config.LoadConfig(cwd, usr, fileAccessor)
	utils.HandleCriticalError(err)

	loggerInstance := logger.NewConsoleLogger(conf.Output)
	runnerInstance := runner.NewRunner(runner.NewRunnerArgs{
		CommandList: []commands.CliCommand{
			initc.NewCommand(flag.ExitOnError),
			handle.NewCommand(flag.ExitOnError, fileAccessor),
		},
		Config:     conf,
		ConfigInfo: configInfo,
		Files:      fileAccessor,
		Logger:     loggerInstance,
		SystemUser: usr,
		Cwd:        cwd,
		Executable: appPath,
	})

	if err = runnerInstance.Run(os.Args[1:]); err != nil {
		loggerInstance.Error(err)
		os.Exit(applicationErrorCode)
	}
}

func panicInterceptor() {
	if err := recover(); err != nil {
		fatal := color.New(color.BgRed, color.FgWhite).SprintFunc()
		log.SetOutput(color.Error)
		log.Println(fatal("Fatal error:"))
		log.Println(fatal(err))
		os.Exit(fatalExitCode)
	}
}
