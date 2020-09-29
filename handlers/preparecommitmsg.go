package handlers

import (
	"fisherman/commands"
	"fisherman/utils"
	"regexp"
)

// PrepareCommitMsgHandler is a execute function for prepare-commit-msg hook
func PrepareCommitMsgHandler(ctx *commands.CommandContext, args []string) error {
	config := ctx.Config.PrepareCommitMsgHook

	if config != nil {
		message, isPresented := getPreparedMessage(config.Message, config.BranchRegExp, "TEST")
		if isPresented {
			err := ctx.Files.Write(args[0], message)
			utils.HandleCriticalError(err)
		}
	}

	return nil
}

func getPreparedMessage(message, regexpString, branch string) (string, bool) {
	if !utils.IsEmpty(message) {
		if !utils.IsEmpty(regexpString) {
			return regexp.MustCompile(regexpString).
				ReplaceAllString(branch, message), true
		}

		return message, true
	}

	return "", false
}
