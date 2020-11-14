package actions

import (
	v "fisherman/internal/validation"
	"fisherman/utils"
)

func PrepareMessage(ctx v.SyncValidationContext, message string) (bool, error) {
	if utils.IsNotEmpty(message) {
		args := ctx.Args()
		err := ctx.Files().Write(args[0], message)

		return false, err
	}

	return true, nil
}
