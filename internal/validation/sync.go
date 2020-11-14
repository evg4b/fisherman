package validation

import (
	"fisherman/internal"

	"github.com/hashicorp/go-multierror"
)

type SyncValidator = func(ctx internal.SyncContext) error

func RunSync(ctx internal.SyncContext, validators []SyncValidator) error {
	var result *multierror.Error

	for _, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
