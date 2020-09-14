package main

import (
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"fisherman/infrastructure/reporter"
	"fisherman/runner"
	"os"
	"os/user"
)

func main() {
	fileAccessor := io.NewFileAccessor()
	log := logger.NewConsoleLooger(logger.OutputConfig{
		LogLevel: logger.Debug,
		Colors:   true,
	})
	usr, err := user.Current()
	handleError(err, log)
	rpt := reporter.NewConsoleReporter()
	r := runner.NewRunner(fileAccessor, usr, rpt)
	handleError(r.Run(os.Args), log)
}

func handleError(err error, log logger.Logger) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
