package shellhandlers

import (
	"fisherman/config/hooks"
	inf "fisherman/infrastructure"
	"fisherman/infrastructure/shell"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func ExecParallel(ctx ContextWithStop, sh inf.Shell, scripts hooks.ScriptsConfig) error {
	chanel := make(chan shell.ExecResult)

	go execInternal(chanel, ctx, sh, scripts)

	var multierr *multierror.Error

	for result := range chanel {
		if result.IsSuccessful() {
			printSuccessful(&result)
		} else {
			printError(&result)
			if !result.IsCanceled() {
				multierr = multierror.Append(multierr, fmt.Errorf("[%s] %s", result.Name, result.Error))
			}
		}
	}

	return multierr.ErrorOrNil()
}
