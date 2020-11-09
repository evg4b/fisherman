package common

import (
	"context"
	"fisherman/config/hooks"
	inf "fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/utils"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type ContextWithStop interface {
	context.Context
	Stop()
}

func ExecCommandsParallel(ctx ContextWithStop, sh inf.Shell, scripts hooks.ScriptsConfig) error {
	chanel := make(chan shell.ExecResult)

	go execInternal(chanel, ctx, sh, scripts)

	var multierr *multierror.Error

	for result := range chanel {
		if result.IsSuccessful() {
			log.Infof("[%s] exit with code %d (executed in %s)", result.Name, result.ExitCode, result.Time)
			if utils.IsNotEmpty(result.Output) {
				log.Info(result.Output)
			}
		} else {
			multierr = multierror.Append(multierr, result.Error)
			log.Errorf("[%s] exit with code %d (executed in %s)", result.Name, result.ExitCode, result.Time)
			if utils.IsNotEmpty(result.Output) {
				log.Error(result.Output)
			}
		}
	}

	return multierr.ErrorOrNil()
}

func execInternal(chanel chan shell.ExecResult, ctx ContextWithStop, sh inf.Shell, scripts hooks.ScriptsConfig) {
	var wg sync.WaitGroup
	for scriptName, shellScript := range scripts {
		wg.Add(1)
		go func(name string, script hooks.ScriptConfig) {
			defer wg.Done()

			result := sh.Exec(ctx, shell.ScriptConfig{
				Name:     name,
				Commands: script.Commands,
				Env:      script.Env,
				Output:   script.Output,
			})

			chanel <- result

			if !result.IsSuccessful() {
				ctx.Stop()
			}
		}(scriptName, shellScript)
	}
	wg.Wait()
	close(chanel)
}
