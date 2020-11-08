// nolint
package handlers

import (
	"fisherman/clicontext"
	"fisherman/handlers/common"
	"fisherman/infrastructure/log"

	"github.com/hashicorp/go-multierror"
	"github.com/mkideal/pkg/errors"
)

// PrePushHandler is a handler for pre-push hook
func PrePushHandler(ctx *clicontext.CommandContext, args []string) error {
	config := ctx.Config.PrePushHook
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
		log.Info(result.Output)
		if result.Err != nil {
			multierr = multierror.Append(multierr, result.Err)
		}
	}

	return multierr.ErrorOrNil()
}
