package clicontext

import (
	"fisherman/config/hooks"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/utils"

	"github.com/imdario/mergo"
	"github.com/mkideal/pkg/errors"
)

func (ctx *CommandContext) Variables() map[string]interface{} {
	ctx.preload()

	return ctx.variables
}

func (ctx *CommandContext) preload() {
	if ctx.variables == nil {
		user, err := ctx.Repository.GetUser()
		utils.HandleCriticalError(err)

		ctx.variables = map[string]interface{}{
			"FishermanVersion": constants.Version,
			"CWD":              ctx.App.Cwd,
			"UserName":         user.UserName,
			"Email":            user.Email,
		}

		err = mergo.Map(&ctx.variables, ctx.globalVariables)
		utils.HandleCriticalError(err)
	}
}

func (ctx *CommandContext) LoadAdditionalVariables(variables *hooks.Variables) error {
	ctx.preload()
	err := ctx.load(ctx.Repository.GetCurrentBranch, variables.GetFromBranch)
	if err != nil {
		return err
	}

	err = ctx.load(ctx.Repository.GetLastTag, variables.GetFromTag)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *CommandContext) load(
	source func() (string, error),
	load func(string) (map[string]interface{}, error),
) error {
	sourceString, err := source()
	if err != nil {
		log.Debugf("Failed getting current branch: %s\n%s", err, errors.Wrap(err))

		return err
	}

	additionalValues, err := load(sourceString)
	if err != nil {
		log.Debugf("Failed getting variables from branch: %s\n%s", err, errors.Wrap(err))

		return err
	}

	err = mergo.MergeWithOverwrite(&ctx.variables, additionalValues)
	if err != nil {
		log.Debugf("Failed merging variables: %s\n%s", err, errors.Wrap(err))

		return err
	}

	return nil
}
