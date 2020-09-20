package handlers

import (
	"fisherman/commands/context"
	"fisherman/config/hooks"
	"fisherman/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// CommitMsgHandler is a handler for commit-msg hook
func CommitMsgHandler(ctx *context.CommandContext, args []string) error {
	config := ctx.Config.Hooks.CommitMsgHook
	if config == nil {
		ctx.Logger.Debug("CommitMsgHook is not presented.")
		return nil
	}

	ctx.Logger.Debug("CommitMsgHook is presented.")
	commitMessage, err := ctx.Files.Read(args[0])
	utils.HandleCriticalError(err)

	if utils.IsEmpty(config.StaticMessage) {
		ctx.Logger.Debug("Static message is presented.")
		err := ctx.Files.Write(args[0], config.StaticMessage)
		utils.HandleCriticalError(err)
		ctx.Logger.Debug("Static message was writted.")
		return nil
	}

	ctx.Logger.Debug("Starting validation.")

	return validateMessage(commitMessage, config).ErrorOrNil()
}

func validateMessage(message string, config *hooks.CommitMsgHookConfig) *multierror.Error {
	var result *multierror.Error

	if config.NotEmpty && utils.IsEmpty(message) {
		err := fmt.Errorf("Commit message should not be empty")
		result = multierror.Append(result, err)
	}

	if !utils.IsEmpty(config.CommitPrefix) && !strings.HasPrefix(message, config.CommitPrefix) {
		err := fmt.Errorf("Commit message should have prefix '%s'", config.CommitPrefix)
		result = multierror.Append(result, err)
	}

	if !utils.IsEmpty(config.CommitSuffix) && !strings.HasSuffix(message, config.CommitSuffix) {
		err := fmt.Errorf("Commit message should have suffix '%s'", config.CommitSuffix)
		result = multierror.Append(result, err)
	}

	if !utils.IsEmpty(config.CommitRegexp) {
		matched, err := regexp.MatchString(config.CommitRegexp, message)
		utils.HandleCriticalError(err)

		if !matched {
			err := fmt.Errorf("Commit message should be matched regular expression '%s'", config.CommitRegexp)
			result = multierror.Append(result, err)
		}
	}

	return result
}
