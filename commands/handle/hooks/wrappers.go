package hooks

import (
	"fisherman/config/hooks"
	"fisherman/infrastructure/shell"
	v "fisherman/internal/validation"
	"fisherman/validators"
)

type validatorWithString = func(v.SyncValidationContext, string) error

func stringWrapper(validator validatorWithString, config string) v.SyncValidator {
	return func(ctx v.SyncValidationContext) error {
		return validator(ctx, config)
	}
}

type validatorWithBool = func(v.SyncValidationContext, bool) error

func boolWrapper(validator validatorWithBool, config bool) v.SyncValidator {
	return func(ctx v.SyncValidationContext) error {
		return validator(ctx, config)
	}
}

func scriptWrapper(scripts hooks.ScriptsConfig) []v.AsyncValidator {
	var validatorList = []v.AsyncValidator{}
	for name, script := range scripts {
		shellScript := shell.ScriptConfig{
			Name:     name,
			Commands: script.Commands,
			Env:      script.Env,
			Output:   true,
		}
		validatorList = append(validatorList, func(ctx v.AsyncValidationContext) v.AsyncValidationResult {
			return validators.ScriptValidator(ctx, shellScript)
		})
	}

	return validatorList
}
