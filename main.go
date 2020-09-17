package main

import (
	"fisherman/config"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"fisherman/runner"
	"fmt"
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
	handleFatalError(err)
	cwd, err := os.Getwd()
	handleFatalError(err)
	conf, err := config.LoadConfig(cwd, usr, fileAccessor)
	handleFatalError(err)
	log := logger.NewConsoleLogger(conf.Config.Output)
	r := runner.NewRunner(fileAccessor, usr, log)

	err = r.Run(conf, os.Args)
	if err != nil {
		log.Error(err)
		os.Exit(applicationErrorCode)
	}
}

func handleFatalError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicInterceptor() {
	if err := recover(); err != nil {
		print := color.New(color.BgRed, color.FgWhite).PrintlnFunc()
		print("Fatal error:")
		print(err)
		os.Exit(fatalExitCode)
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
	}
}
