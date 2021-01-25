package hookfactory

import (
	hooks "fisherman/configuration"
	"fisherman/infrastructure/shell"
	"fisherman/internal"
	"fisherman/internal/expression"
	"fisherman/internal/validation"
	"fisherman/utils"
	"fisherman/validators"
)

func scriptWrapper(scripts hooks.ScriptsConfig, engine expression.Engine) []validation.AsyncValidator {
	var validatorList = []validation.AsyncValidator{}
	for name, script := range scripts {
		if !utils.IsEmpty(script.Condition) {
			condition, err := engine.Eval(script.Condition)
			utils.HandleCriticalError(err)
			if !condition {
				continue
			}
		}

		validatorList = append(validatorList, func(ctx internal.ExecutionContext) validation.AsyncValidationResult {
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
