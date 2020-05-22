// nolint: errcheck, dupl
package shell_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/evg4b/fisherman/pkg/guards"
	"github.com/evg4b/fisherman/testing/mocks"
	"os/exec"
	"strings"
	"testing"

	. "github.com/evg4b/fisherman/pkg/shell"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func TestHost_Start(t *testing.T) {
	t.Run("second call start returns error", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		guards.NoError(err)
		defer host.Terminate()

		err = host.Start()

		require.EqualError(t, err, "host already started")
	})

	t.Run("start returns error after print script", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		defer host.Terminate()

		fmt.Fprintln(host, "echo 1")
		err := host.Start()

		require.EqualError(t, err, "host already started")
	})

	t.Run("stdin pipe error", func(t *testing.T) {
		command := exec.Command("echo", "test")
		command.Stdin = strings.NewReader("test")
		strategy := mocks.NewShellStrategyMock(t).
			GetCommandMock.Return(command).
			GetEncodingMock.Return(encoding.Nop)
		host := NewHost(context.TODO(), strategy)
		defer host.Terminate()

		err := host.Start()

		require.EqualError(t, err, "exec: Stdin already set")
	})
}

func TestHost_Run(t *testing.T) {
	t.Run("return error for empty string", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		defer host.Terminate()

		err := host.Run("")

		require.EqualError(t, err, "script can not be empty")
	})

	t.Run("return error for started host", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		guards.NoError(err)
		defer host.Terminate()

		err = host.Run("echo 'test'")

		require.EqualError(t, err, "host already started")
	})
}

func TestHost_Write(t *testing.T) {
	t.Run("write with fmt correctly", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		defer host.Terminate()

		_, err := fmt.Fprintln(host, "echo 'test'")

		require.NoError(t, err)
	})

	t.Run("stdin pipe error", func(t *testing.T) {
		command := exec.Command("echo", "test")
		command.Stdin = strings.NewReader("test")
		strategy := mocks.NewShellStrategyMock(t).
			GetCommandMock.Return(command).
			GetEncodingMock.Return(encoding.Nop)
		host := NewHost(context.TODO(), strategy)
		defer host.Terminate()

		_, err := fmt.Fprintln(host, "echo 'test'")

		require.EqualError(t, err, "exec: Stdin already set")
	})

	t.Run("write in call encoder", func(t *testing.T) {
		encoderTransformer := mocks.NewTransformerMock(t).
			ResetMock.Return().
			TransformMock.Set(transform.Nop.Transform)

		decoderTransformer := mocks.NewTransformerMock(t).
			ResetMock.Return().
			TransformMock.Set(transform.Nop.Transform)

		encodingMock := mocks.NewEncodingMock(t).
			NewDecoderMock.Return(&encoding.Decoder{Transformer: decoderTransformer}).
			NewEncoderMock.Return(&encoding.Encoder{Transformer: encoderTransformer})

		host := NewHost(context.TODO(), Default(), WithEncoding(encodingMock))
		defer host.Terminate()

		_, err := fmt.Fprintln(host, "echo 'test'")

		_ = host.Wait()

		require.NoError(t, err)
		assert.NotEmpty(t, encoderTransformer.TransformMock.Calls())
		assert.Empty(t, decoderTransformer.TransformMock.Calls())
	})

	t.Run("write in call encoder and decoder", func(t *testing.T) {
		encoderTransformer := mocks.NewTransformerMock(t).
			ResetMock.Set(transform.Nop.Reset).
			TransformMock.Set(transform.Nop.Transform)

		decoderTransformer := mocks.NewTransformerMock(t).
			ResetMock.Set(transform.Nop.Reset).
			TransformMock.Set(transform.Nop.Transform)

		encodingMock := mocks.NewEncodingMock(t).
			NewDecoderMock.Return(&encoding.Decoder{Transformer: decoderTransformer}).
			NewEncoderMock.Return(&encoding.Encoder{Transformer: encoderTransformer})

		buff := &bytes.Buffer{}
		host := NewHost(context.TODO(), Default(), WithEncoding(encodingMock), WithStdout(buff), WithStderr(buff))
		defer host.Terminate()

		_, err := fmt.Fprintln(host, "echo 'test'")

		_ = host.Wait()

		require.NoError(t, err)
		assert.NotEmpty(t, encoderTransformer.TransformMock.Calls())
		assert.NotEmpty(t, decoderTransformer.TransformMock.Calls())
	})
}

func TestHost_Close(t *testing.T) {
	t.Run("not started host", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())

		err := host.Close()

		require.NoError(t, err)
	})

	t.Run("correctly closed stdin", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		guards.NoError(err)
		defer host.Terminate()

		err = host.Close()

		require.NoError(t, err)
	})

	t.Run("called multiple times", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		err := host.Start()
		guards.NoError(err)
		defer host.Terminate()

		assert.NotPanics(t, func() {
			err = host.Close()
			require.NoError(t, err)

			err = host.Close()
			require.NoError(t, err)
		})
	})

	t.Run("encoding error", func(t *testing.T) {
		transformer := mocks.NewTransformerMock(t).
			ResetMock.Set(transform.Nop.Reset).
			TransformMock.Set(func(dst, src []byte, atEOF bool) (nDst int, nSrc int, err error) {
			return 0, 0, errors.New("encoding error")
		})

		encodingMock := mocks.NewEncodingMock(t).
			NewDecoderMock.Return(&encoding.Decoder{Transformer: transformer}).
			NewEncoderMock.Return(&encoding.Encoder{Transformer: transformer})

		host := NewHost(context.TODO(), Default(), WithEncoding(encodingMock))
		defer host.Terminate()

		_ = host.Start()

		actual := host.Close()

		require.EqualError(t, actual, "encoding error")
		assert.NotEmpty(t, transformer.TransformMock.Calls())
	})
}

func TestHost_Terminate(t *testing.T) {
	t.Run("returns errors where host not started", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())

		err := host.Terminate()

		require.EqualError(t, err, "host is not started")
	})

	t.Run("terminate host correctly", func(t *testing.T) {
		host := NewHost(context.TODO(), Default())
		_ = host.Start()

		err := host.Terminate()

		require.NoError(t, err)
	})
}

func TestHost_Wait(t *testing.T) {
	t.Run("return encoding error", func(t *testing.T) {
		transformer := mocks.NewTransformerMock(t).
			ResetMock.Set(transform.Nop.Reset).
			TransformMock.Set(func(dst, src []byte, atEOF bool) (nDst int, nSrc int, err error) {
			return 0, 0, errors.New("encodingMock error")
		})

		encodingMock := mocks.NewEncodingMock(t).
			NewDecoderMock.Return(&encoding.Decoder{Transformer: transformer}).
			NewEncoderMock.Return(&encoding.Encoder{Transformer: transformer})

		host := NewHost(context.TODO(), Default(), WithEncoding(encodingMock))
		defer host.Terminate()
		_ = host.Start()

		actual := host.Wait()

		require.EqualError(t, actual, "encodingMock error")
		assert.NotEmpty(t, transformer.TransformMock.Calls())
	})
}
