package validators

import (
	"fisherman/infrastructure/shell"
	"fisherman/internal"
	"fisherman/internal/validation"
)

func ScriptValidator(ctx internal.AsyncContext, script shell.ShScriptConfig) validation.AsyncValidationResult {
	sh := ctx.Shell()
	result := sh.Exec(ctx, script)

	return validation.AsyncValidationResult{
		Name:  result.Name,
		Error: result.Error,
		Time:  result.Time,
	}
}
