// nolint
package handlers

import (
	"fisherman/clicontext"
	"fisherman/handlers/common"
	"fisherman/infrastructure/log"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/mkideal/pkg/errors"
)

// PreCommitHandler is a handler for pre-commit hook
func PreCommitHandler(ctx *clicontext.CommandContext, args []string) error {
	config := ctx.Config.PreCommitHook
	err := ctx.LoadAdditionalVariables(&config.Variables)
	if err != nil {
		log.Debugf("Additional variables loading filed: %s\n%s", err, errors.Wrap(err))

		return err
	}

	config.Compile(ctx.Variables())

	var multierr *multierror.Error
	results := common.ExecCommandsParallel(ctx.Shell, config.Shell)
	for key, result := range results {
		log.Infof("[%s] exited with code %d", key, result.ExitCode)
		if len(result.Output) > 0 {
			log.Info(result.Output)
		}

		if result.Err != nil {
			multierr = multierror.Append(multierr, result.Err)
		}

		if result.ExitCode != 0 {
			err = fmt.Errorf("script %s exited with code %d", key, result.ExitCode)
			multierr = multierror.Append(multierr, err)
		}
	}

	return multierr.ErrorOrNil()
}
