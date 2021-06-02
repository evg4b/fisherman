package utils_test

import (
	"errors"
	"fisherman/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecWithTime(t *testing.T) {
	duration, err := utils.ExecWithTime(func() error {
		time.Sleep(time.Millisecond * 1)

		return nil
	})

	assert.NoError(t, err)
	assert.Greater(t, int(duration), 0)
}

func TestExecWithTime_Error(t *testing.T) {
	duration, err := utils.ExecWithTime(func() error {
		time.Sleep(time.Millisecond * 1)

		return errors.New("TestError")
	})

	assert.Error(t, err, "TestError")
	assert.Greater(t, int(duration), 0)
}
