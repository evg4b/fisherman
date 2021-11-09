package shell

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: remove this test file

func TestExecWithTime(t *testing.T) {
	duration, err := execWithTime(func() error {
		time.Sleep(time.Millisecond * 1)

		return nil
	})

	assert.NoError(t, err)
	assert.Greater(t, int(duration), 0)
}

func TestExecWithTime_Error(t *testing.T) {
	duration, err := execWithTime(func() error {
		time.Sleep(time.Millisecond * 1)

		return errors.New("TestError")
	})

	assert.EqualError(t, err, "TestError")
	assert.Greater(t, int(duration), 0)
}
