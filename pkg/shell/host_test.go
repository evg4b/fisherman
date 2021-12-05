package shell_test

import (
	"context"
	. "fisherman/pkg/shell"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHost_Start(t *testing.T) {
	t.Run("second call start returns error", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		defer host.Terminate() // nolint: errcheck
		assert.NoError(t, err)

		err = host.Start()

		assert.EqualError(t, err, "shell host: already started")
	})

	t.Run("start returns error after print script", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		defer host.Terminate() // nolint: errcheck

		fmt.Fprintln(host, "echo 1")
		err := host.Start()
		assert.EqualError(t, err, "shell host: already started")
	})
}
