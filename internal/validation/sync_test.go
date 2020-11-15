package validation_test

import (
	"errors"
	"fisherman/internal"
	"fisherman/internal/validation"
	"fisherman/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSync(t *testing.T) {
	callCont := 0

	validators := []validation.SyncValidator{
		func(ctx internal.SyncContext) error {
			callCont++

			return errors.New("test error")
		},
		func(ctx internal.SyncContext) error {
			callCont++

			return nil
		},
	}

	ctx := mocks.NewSyncContextMock(t)

	err := validation.RunSync(ctx, validators)

	assert.EqualError(t, err, "1 error occurred:\n\t* test error\n\n")
}
