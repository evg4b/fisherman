package shellhandlers

import (
	"context"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/utils"
)

type ContextWithStop interface {
	context.Context
	Stop()
}

func printError(result *shell.ExecResult) {
	if !result.IsCanceled() {
		log.Errorf("[%s] exit with code %d (executed in %s)", result.Name, result.ExitCode, result.Time)
		if utils.IsNotEmpty(result.Output) {
			log.Error(result.Output)
		}
	} else {
		log.Infof("[%s] was skipped", result.Name)
	}
}

func printSuccessful(result *shell.ExecResult) {
	log.Infof("[%s] exit with code %d (executed in %s)", result.Name, result.ExitCode, result.Time)
	if utils.IsNotEmpty(result.Output) {
		log.Info(result.Output)
	}
}
