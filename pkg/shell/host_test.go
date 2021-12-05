// nolint: errcheck
package shell_test

import (
	"context"
	"fisherman/pkg/guards"
	. "fisherman/pkg/shell"
	"fisherman/testing/mocks"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHost_Start(t *testing.T) {
	t.Run("second call start returns error", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		guards.NoError(err)
		defer host.Terminate()

		err = host.Start()

		assert.EqualError(t, err, "host already started")
	})

	t.Run("start returns error after print script", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		defer host.Terminate()

		fmt.Fprintln(host, "echo 1")
		err := host.Start()

		assert.EqualError(t, err, "host already started")
	})

	t.Run("stdin pipe error", func(t *testing.T) {
		command := exec.Command("echo", "test")
		command.Stdin = strings.NewReader("test")
		strategy := mocks.NewShellStrategyMock(t).
			GetCommandMock.Return(command)
		host := NewHost(context.TODO(), strategy)
		defer host.Terminate()

		err := host.Start()

		assert.EqualError(t, err, "exec: Stdin already set")
	})
}

func TestHost_Run(t *testing.T) {
	t.Run("return error for empty string", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		defer host.Terminate()

		err := host.Run("")

		assert.EqualError(t, err, "script can not be empty")
	})

	t.Run("return error for started host", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		guards.NoError(err)
		defer host.Terminate()

		err = host.Run("echo 'test'")

		assert.EqualError(t, err, "host already started")
	})
}

func TestHost_Write(t *testing.T) {
	t.Run("write with fmt correctly", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		defer host.Terminate()

		_, err := fmt.Fprintln(host, "echo 'test'")

		assert.NoError(t, err)
	})

	t.Run("stdin pipe error", func(t *testing.T) {
		command := exec.Command("echo", "test")
		command.Stdin = strings.NewReader("test")
		strategy := mocks.NewShellStrategyMock(t).
			GetCommandMock.Return(command)
		host := NewHost(context.TODO(), strategy)
		defer host.Terminate()

		_, err := fmt.Fprintln(host, "echo 'test'")

		assert.EqualError(t, err, "exec: Stdin already set")
	})
}

func TestHost_Close(t *testing.T) {
	t.Run("not started host", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())

		err := host.Close()

		assert.NoError(t, err)
	})

	t.Run("correctly closed stdin", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		guards.NoError(err)
		defer host.Terminate()

		err = host.Close()

		assert.NoError(t, err)
	})

	t.Run("called multiple times", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		guards.NoError(err)
		defer host.Terminate()

		assert.NotPanics(t, func() {
			err = host.Close()
			assert.NoError(t, err)

			err = host.Close()
			assert.NoError(t, err)
		})
	})
}

func TestHost_Terminate(t *testing.T) {
	t.Run("returns errors where host not started", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())

		err := host.Terminate()

		assert.EqualError(t, err, "host is not started")
	})

	t.Run("terminate host correctly", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		_ = host.Start()

		err := host.Terminate()

		assert.NoError(t, err)
	})
}
