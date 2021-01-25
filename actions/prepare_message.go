package actions

import (
	"fisherman/internal"
	"fisherman/utils"
)

func PrepareMessage(ctx internal.ExecutionContext, message string) (bool, error) {
	if !utils.IsEmpty(message) {
		args := ctx.Args()
		files := ctx.Files()

		return false, files.Write(args[0], message)
	}

	return true, nil
}
