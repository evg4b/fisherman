package handlers

import (
	"fisherman/commands/context"
	"fisherman/config/hooks"
	"fisherman/utils"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// CommitMsgHandler is a handler for commit-msg hook
func CommitMsgHandler(ctx *context.CommandContext, args []string) error {
	appConfig := ctx.GetHookConfiguration()
	config := appConfig.CommitMsgHook
	if config == nil {
		ctx.Logger.Debug("CommitMsgHook is not presented.")
		return nil
	}

	ctx.Logger.Debug("CommitMsgHook is presented.")
	commitMessage, err := ctx.FileAccessor.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("Can't open '%s' file: %e", args[0], err)
	}

	if utils.IsEmpty(config.StaticMessage) {
		ctx.Logger.Debug("Static message is presented.")
		err := ctx.FileAccessor.WriteFile(args[0], config.StaticMessage)
		if err != nil {
			return fmt.Errorf("Can't write '%s' file: %e", args[0], err)
		}

		ctx.Logger.Debug("Static message was writted.")
		return nil
	}

	ctx.Logger.Debug("Starting validation.")

	return validateMessage(commitMessage, config).ErrorOrNil()
}

func validateMessage(message string, config *hooks.CommitMsgHookConfig) *multierror.Error {
	var result *multierror.Error

	if config.NotEmpty && utils.IsEmpty(message) {
		err := fmt.Errorf("Commit comment should not be empty")
		result = multierror.Append(result, err)
	}

	if !utils.IsEmpty(config.CommitPrefix) && !strings.HasPrefix(message, config.CommitPrefix) {
		err := fmt.Errorf("Commit should have prefix '%s'", config.CommitPrefix)
		result = multierror.Append(result, err)
	}

	if !utils.IsEmpty(config.CommitSuffix) && !strings.HasSuffix(message, config.CommitSuffix) {
		err := fmt.Errorf("Commit should have suffix '%s'", config.CommitSuffix)
		result = multierror.Append(result, err)
	}

	return result
}
