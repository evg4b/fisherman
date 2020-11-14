package validators

import (
	"fisherman/infrastructure/shell"
	v "fisherman/internal/validation"
)

func ScriptValidator(ctx v.AsyncValidationContext, script shell.ScriptConfig) v.AsyncValidationResult {
	sh := ctx.Shell()
	result := sh.Exec(ctx, script)

	return v.AsyncValidationResult{
		Name:  result.Name,
		Error: result.Error,
		Time:  result.Time,
	}
}
