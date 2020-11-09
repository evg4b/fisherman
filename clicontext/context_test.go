package clicontext_test

import (
	"context"
	"fisherman/clicontext"
	"fisherman/config"
	ctx_mock "fisherman/mocks/pkg/context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommandContext_Deadline(t *testing.T) {
	expectedDuration := time.Time{}
	expectedResult := true

	fakeCtx := getFakeCtx(expectedDuration, expectedResult)

	ctx := clicontext.NewContext(fakeCtx, clicontext.Args{
		Config: &config.FishermanConfig{},
	})

	duration, ok := ctx.Deadline()

	assert.Equal(t, expectedDuration, duration)
	assert.Equal(t, expectedResult, ok)
}

func TestCommandContext_Done(t *testing.T) {
	fakeCtx := getFakeCtx(time.Time{}, false)

	ctx := clicontext.NewContext(fakeCtx, clicontext.Args{
		Config: &config.FishermanConfig{},
	})

	chanel := ctx.Done()

	fakeCtx.AssertCalled(t, "Done")
	assert.NotNil(t, chanel)
}

func TestCommandContext_Err(t *testing.T) {
	baseCtx, cancel := context.WithCancel(context.TODO())
	cancel()

	ctx := clicontext.NewContext(baseCtx, clicontext.Args{
		Config: &config.FishermanConfig{},
	})

	err := ctx.Err()

	assert.Equal(t, context.Canceled, err)
}

func TestCommandContext_Value(t *testing.T) {
	key := "testKey"
	expectedValue := "testValue"

	// nolint: go-lint
	baseCtx := context.WithValue(context.Background(), key, expectedValue)

	ctx := clicontext.NewContext(baseCtx, clicontext.Args{
		Config: &config.FishermanConfig{},
	})

	value := ctx.Value(key)

	assert.Equal(t, expectedValue, value)
}

func TestCommandContext_Stop(t *testing.T) {
	ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
		Config: &config.FishermanConfig{},
	})

	ctx.Stop()

	assert.Equal(t, context.Canceled, ctx.Err())
}

func getFakeCtx(expectedTime time.Time, ok bool) *ctx_mock.Context {
	fakeCtx := ctx_mock.Context{}
	fakeCtx.On("Deadline").Return(expectedTime, ok)
	fakeCtx.On("Done").Return(make(<-chan struct{}))
	fakeCtx.On("Value", mock.Anything).Return("test")

	return &fakeCtx
}
