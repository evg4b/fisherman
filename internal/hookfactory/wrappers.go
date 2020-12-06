package hookfactory

import (
	hooks "fisherman/configuration"
	"fisherman/infrastructure/shell"
	"fisherman/internal"
	"fisherman/internal/validation"
	"fisherman/validators"
)

type stringF = func(internal.SyncContext, string) error

func stringWrapper(validator stringF, config string) validation.SyncValidator {
	return func(ctx internal.SyncContext) error {
		return validator(ctx, config)
	}
}

type boolF = func(internal.SyncContext, bool) error

func boolWrapper(validator boolF, config bool) validation.SyncValidator {
	return func(ctx internal.SyncContext) error {
		return validator(ctx, config)
	}
}

func scriptWrapper(scripts hooks.ScriptsConfig) []validation.AsyncValidator {
	var validatorList = []validation.AsyncValidator{}
	for name, script := range scripts {
		validatorList = append(validatorList, func(ctx internal.AsyncContext) validation.AsyncValidationResult {
			return validators.ScriptValidator(ctx, script.Shell, shell.ShScriptConfig{
				Name:     name,
				Commands: script.Commands,
				Env:      script.Env,
				Output:   true,
				Dir:      script.Dir,
			})
		})
	}

	return validatorList
}
