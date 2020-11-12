package handlers

import (
	"errors"
	"fisherman/clicontext"
	"fisherman/config/hooks"
	"fisherman/infrastructure/log"
	"fisherman/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// CommitMsgHandler is a handler for commit-msg hook
func CommitMsgHandler(ctx *clicontext.CommandContext, args []string) error {
	if len(args) < 1 {
		return errors.New("commit message file argument is not presented")
	}

	config := &ctx.Config.CommitMsgHook

	err := ctx.LoadAdditionalVariables(&config.Variables)
	if err != nil {
		log.Debugf("Additional variables loading filed: %s\n%s", err, err)

		return err
	}
	log.Debug("Additional variables was loaded")

	utils.FillTemplate(&config.MessagePrefix, ctx.Variables())
	utils.FillTemplate(&config.MessageSuffix, ctx.Variables())
	utils.FillTemplate(&config.StaticMessage, ctx.Variables())
	log.Debug("Templates was compiled")

	log.Debugf("Reading commit message file %s", args[0])
	commitMessage, err := ctx.Files.Read(args[0])
	utils.HandleCriticalError(err)
	log.Debugf("Commit message file was successful read")

	if utils.IsNotEmpty(config.StaticMessage) {
		log.Debug("Static message is presented.")
		err := ctx.Files.Write(args[0], config.StaticMessage)
		utils.HandleCriticalError(err)
		log.Debug("Static message was writted.")

		return nil
	}

	log.Debug("Static message is not presented. Starting validation.")

	return validateMessage(strings.TrimSpace(commitMessage), config).ErrorOrNil()
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
