// nolint
package handlers

import (
	"fisherman/clicontext"
	"fisherman/handlers/common"
	"fisherman/infrastructure/log"

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

	return common.ExecCommandsParallel(ctx, ctx.Shell, config.Shell)
}
