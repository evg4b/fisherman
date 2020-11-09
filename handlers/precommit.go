package handlers

import (
	"fisherman/clicontext"
	"fisherman/handlers/shellhandlers"
	"fisherman/infrastructure/log"

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

	return shellhandlers.ExecParallel(ctx, ctx.Shell, config.Shell)
}
