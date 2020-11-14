package validators

import (
	"fisherman/infrastructure/shell"
	i "fisherman/internal"
	v "fisherman/internal/validation"
)

func ScriptValidator(ctx i.AsyncContext, script shell.ScriptConfig) v.AsyncValidationResult {
	sh := ctx.Shell()
	result := sh.Exec(ctx, script)

	return v.AsyncValidationResult{
		Name:  result.Name,
		Error: result.Error,
		Time:  result.Time,
	}
}
