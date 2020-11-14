package actions

import (
	"fisherman/internal"
	"fisherman/utils"
)

func PrepareMessage(ctx internal.SyncContext, message string) (bool, error) {
	if utils.IsNotEmpty(message) {
		args := ctx.Args()
		err := ctx.Files().Write(args[0], message)

		return false, err
	}

	return true, nil
}
