package handlers

import (
	"fisherman/commands"
	"fisherman/config/hooks"
	"fisherman/infrastructure/logger"
	"fisherman/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// CommitMsgHandler is a handler for commit-msg hook
func CommitMsgHandler(ctx *commands.CommandContext, args []string) error {
	config := &ctx.Config.CommitMsgHook

	err := ctx.LoadAdditionalVariables(&config.Variables)
	if err != nil {
		return err
	}

	utils.FillTemplate(&config.MessagePrefix, ctx.Variables)
	utils.FillTemplate(&config.MessageSuffix, ctx.Variables)
	utils.FillTemplate(&config.StaticMessage, ctx.Variables)

	logger.Debug("CommitMsgHook is presented.")
	commitMessage, err := ctx.Files.Read(args[0])
	utils.HandleCriticalError(err)

	if utils.IsNotEmpty(config.StaticMessage) {
		logger.Debug("Static message is presented.")
		err := ctx.Files.Write(args[0], config.StaticMessage)
		utils.HandleCriticalError(err)
		logger.Debug("Static message was writted.")

		return nil
	}

	logger.Debug("Starting validation.")

	return validateMessage(commitMessage, config).ErrorOrNil()
}

func validateMessage(message string, config *hooks.CommitMsgHookConfig) *multierror.Error {
	var result *multierror.Error

	if config.NotEmpty && utils.IsEmpty(message) {
		err := fmt.Errorf("commit message should not be empty")
		result = multierror.Append(result, err)
	}

	if utils.IsNotEmpty(config.MessagePrefix) && !strings.HasPrefix(message, config.MessagePrefix) {
		err := fmt.Errorf("commit message should have prefix '%s'", config.MessagePrefix)
		result = multierror.Append(result, err)
	}

	if utils.IsNotEmpty(config.MessageSuffix) && !strings.HasSuffix(message, config.MessageSuffix) {
		err := fmt.Errorf("commit message should have suffix '%s'", config.MessageSuffix)
		result = multierror.Append(result, err)
	}

	if utils.IsNotEmpty(config.MessageRegexp) {
		matched, err := regexp.MatchString(config.MessageRegexp, message)
		utils.HandleCriticalError(err)

		if !matched {
			err := fmt.Errorf("commit message should be matched regular expression '%s'", config.MessageRegexp)
			result = multierror.Append(result, err)
		}
	}

	return result
}
