package validation_test

import (
	"context"
	"errors"
	"fisherman/internal"
	"fisherman/internal/validation"
	"fisherman/mocks"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestRunAsync_Empty(t *testing.T) {
	err := validation.RunAsync(mocks.NewAsyncContextMock(t), []validation.AsyncValidator{})

	assert.NoError(t, err)
}

func TestRunAsync(t *testing.T) {
	ctx := mocks.NewAsyncContextMock(t).ErrMock.Return(nil).StopMock.Return()

	validators := []validation.AsyncValidator{
		func(ctx internal.AsyncContext) validation.AsyncValidationResult {
			return validation.AsyncValidationResult{
				Name:  "test",
				Error: nil,
				Time:  time.Hour,
			}
		},
		func(ctx internal.AsyncContext) validation.AsyncValidationResult {
			return validation.AsyncValidationResult{
				Name:  "test-2",
				Error: errors.New("error-1"),
				Time:  time.Hour,
			}
		},
		func(ctx internal.AsyncContext) validation.AsyncValidationResult {
			return validation.AsyncValidationResult{
				Name:  "test-3",
				Error: nil,
				Time:  time.Hour,
			}
		},
	}

	err := validation.RunAsync(ctx, validators)

	assert.EqualValues(
		t,
		err.(*multierror.Error).Errors,
		[]error{errors.New("[test-2] error-1")},
	)
}

func TestRunAsync_Canceled(t *testing.T) {
	ctx := mocks.NewAsyncContextMock(t).
		ErrMock.Return(context.Canceled).
		StopMock.Return()

	validators := []validation.AsyncValidator{
		func(ctx internal.AsyncContext) validation.AsyncValidationResult {
			return validation.AsyncValidationResult{
				Name:  "test-2",
				Error: nil,
				Time:  time.Hour,
			}
		},
	}

	err := validation.RunAsync(ctx, validators)

	assert.NoError(t, err)
}

func TestRunAsync_DeadlineExceeded(t *testing.T) {
	ctx := mocks.NewAsyncContextMock(t).
		ErrMock.Return(context.DeadlineExceeded).
		StopMock.Return()

	validators := []validation.AsyncValidator{
		func(ctx internal.AsyncContext) validation.AsyncValidationResult {
			return validation.AsyncValidationResult{
				Name:  "test-2",
				Error: nil,
				Time:  time.Hour,
			}
		},
	}

	err := validation.RunAsync(ctx, validators)

	assert.EqualError(t, err, "1 error occurred:\n\t* [test-2] context deadline exceeded\n\n")
}
