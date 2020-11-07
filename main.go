package main

import (
	"fisherman/runner"
	"fisherman/utils"
	"os"
)

const fatalExitCode = 1

func main() {
	defer utils.PanicInterceptor(os.Exit, fatalExitCode)
	args, err := runner.GetRunnerArgs()
	utils.HandleCriticalError(err)
	instance := runner.NewRunner(args)
	if err = instance.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
