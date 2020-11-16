package validators

import (
	"fisherman/infrastructure/shell"
	"fisherman/internal"
	"fisherman/internal/validation"
)

func ScriptValidator(
	ctx internal.AsyncContext,
	shellName string,
	script shell.ShScriptConfig,
) validation.AsyncValidationResult {
	sh := ctx.Shell()
	result := sh.Exec(ctx, shellName, script)

	return validation.AsyncValidationResult{
		Name:  result.Name,
		Error: result.Error,
		Time:  result.Time,
	}
}
