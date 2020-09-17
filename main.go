package main

import (
	"fisherman/config"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"fisherman/runner"
	"fisherman/utils"
	"os"
	"os/user"

	"github.com/fatih/color"
)

const fatalExitCode = 1
const applicationErrorCode = 200

func main() {
	defer panicInterceptor()
	fileAccessor := io.NewFileAccessor()
	usr, err := user.Current()
	utils.HandleCriticalError(err)

	cwd, err := os.Getwd()
	utils.HandleCriticalError(err)

	conf, err := config.LoadConfig(cwd, usr, fileAccessor)
	utils.HandleCriticalError(err)

	log := logger.NewConsoleLogger(conf.Config.Output)
	r := runner.NewRunner(fileAccessor, usr, log)

	if err = r.Run(conf, os.Args); err != nil {
		log.Error(err)
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
