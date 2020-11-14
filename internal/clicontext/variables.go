package clicontext

import (
	"fisherman/config/hooks"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/utils"

	"github.com/imdario/mergo"
)

type sourceLoader = func() (string, error)
type variablesLoader = func(string) (map[string]interface{}, error)

// Variables initializes (when variables was not initialized) and returns registered variables in context
func (ctx *CommandContext) Variables() map[string]interface{} {
	ctx.preload()

	return ctx.variables
}

// Variables initializes (when variables was not initialized) and load additional variables in contxt
func (ctx *CommandContext) LoadAdditionalVariables(variables hooks.VariablesExtractor) error {
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

func (ctx *CommandContext) preload() {
	if ctx.variables == nil {
		user, err := ctx.Repository.GetUser()
		utils.HandleCriticalError(err)

		ctx.variables = map[string]interface{}{
			constants.FishermanVersionVariable: constants.Version,
			constants.CwdVariable:              ctx.App.Cwd,
			constants.UserNameVariable:         user.UserName,
			constants.EmailVariable:            user.Email,
		}

		err = mergo.Map(&ctx.variables, ctx.globalVariables)
		utils.HandleCriticalError(err)
	}
}

func (ctx *CommandContext) load(source sourceLoader, load variablesLoader) error {
	sourceString, err := source()
	if err != nil {
		log.Debugf("Failed getting source string: %s", err)

		return err
	}

	additionalValues, err := load(sourceString)
	if err != nil {
		log.Debugf("Failed getting variables from string '%s': %s", sourceString, err)

		return err
	}

	return mergo.MergeWithOverwrite(&ctx.variables, additionalValues)
}
