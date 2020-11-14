package hooks

import (
	"fisherman/config/hooks"
	"fisherman/infrastructure/shell"
	"fisherman/internal"
	v "fisherman/internal/validation"
	"fisherman/validators"
)

type validatorWithString = func(internal.SyncContext, string) error

func stringWrapper(validator validatorWithString, config string) v.SyncValidator {
	return func(ctx internal.SyncContext) error {
		return validator(ctx, config)
	}
}

type validatorWithBool = func(internal.SyncContext, bool) error

func boolWrapper(validator validatorWithBool, config bool) v.SyncValidator {
	return func(ctx internal.SyncContext) error {
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
		validatorList = append(validatorList, func(ctx internal.AsyncContext) v.AsyncValidationResult {
			return validators.ScriptValidator(ctx, shellScript)
		})
	}

	return validatorList
}
