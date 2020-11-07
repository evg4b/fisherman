package clicontext

import (
	"fisherman/config/hooks"
	"fisherman/infrastructure/log"

	"github.com/imdario/mergo"
	"github.com/mkideal/pkg/errors"
)

func (ctx *CommandContext) LoadAdditionalVariables(variables *hooks.Variables) error {
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

	err = mergo.MergeWithOverwrite(&ctx.Variables, additionalValues)
	if err != nil {
		log.Debugf("Failed merging variables: %s\n%s", err, errors.Wrap(err))

		return err
	}

	return nil
}
