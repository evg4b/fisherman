package validation

import (
	"fisherman/infrastructure"
	"io"

	"github.com/hashicorp/go-multierror"
)

type SyncValidationContext interface {
	Files() infrastructure.FileSystem
	Shell() infrastructure.Shell
	Repository() infrastructure.Repository
	Args() []string
	Output() io.Writer
	Message() string
}

type SyncValidator = func(ctx SyncValidationContext) error

func RunSync(ctx SyncValidationContext, validators []SyncValidator) error {
	var result *multierror.Error

	for _, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
