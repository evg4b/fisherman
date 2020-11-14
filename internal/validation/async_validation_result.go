package validation

import (
	"context"
	"time"
)

type AsyncValidationResult struct {
	Name  string
	Error error
	Time  time.Duration
}

func (result *AsyncValidationResult) IsSuccessful() bool {
	return result.Error == nil
}

func (result *AsyncValidationResult) IsCanceled() bool {
	return result.Error == context.Canceled
}
