package sysshell

import (
	"context"
	"fisherman/config/hooks"
	inf "fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type ContextWithStop interface {
	context.Context
	Stop()
}

func ExecParallel(ctx ContextWithStop, sh inf.Shell, scripts hooks.ScriptsConfig) error {
	chanel := make(chan shell.ExecResult)

	go execInternal(chanel, ctx, sh, scripts)

	var multierr *multierror.Error

	for result := range chanel {
		if result.IsCanceled() {
			log.Infof("[%s] was skipped", result.Name)

			continue
		}

		if result.IsSuccessful() {
			log.Infof("[%s] exit with code %d (executed in %s)", result.Name, result.ExitCode, result.Time)
		} else {
			log.Infof("[%s] exit with code %d (executed in %s)", result.Name, result.ExitCode, result.Time)
			multierr = multierror.Append(multierr, fmt.Errorf("[%s] %s", result.Name, result.Error))
		}
	}

	return multierr.ErrorOrNil()
}
