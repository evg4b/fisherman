package clicontext_test

import (
	"context"
	"fisherman/config"
	"fisherman/internal/clicontext"
	"fisherman/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandContext_Deadline(t *testing.T) {
	expectedDuration := time.Time{}
	expectedResult := true

	fakeCtx := getFakeCtx(t, expectedDuration, expectedResult)

	ctx := clicontext.NewContext(fakeCtx, clicontext.Args{
		Config: &config.FishermanConfig{},
	})

	duration, ok := ctx.Deadline()

	assert.Equal(t, expectedDuration, duration)
	assert.Equal(t, expectedResult, ok)
}

func TestCommandContext_Done(t *testing.T) {
	fakeCtx := getFakeCtx(t, time.Time{}, false)

	ctx := clicontext.NewContext(fakeCtx, clicontext.Args{
		Config: &config.FishermanConfig{},
	})

	chanel := ctx.Done()

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

	// nolint: golint, staticcheck
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

func getFakeCtx(t *testing.T, expectedTime time.Time, ok bool) *mocks.ContextWithStopMock {
	return mocks.NewContextWithStopMock(t).
		DeadlineMock.Return(expectedTime, ok).
		DoneMock.Return(make(<-chan struct{})).
		ValueMock.Return("test")
}
