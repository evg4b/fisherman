package internal_test

import (
	"bytes"
	"context"
	"errors"
	"fisherman/internal"
	"fisherman/mocks"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_Files(t *testing.T) {
	fs := mocks.NewFileSystemMock(t)
	ctx := internal.NewInternalContext(
		context.TODO(),
		fs,
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{},
		ioutil.Discard,
	)

	actualFs := ctx.Files()

	assert.Equal(t, fs, actualFs)
}

func TestContext_Shell(t *testing.T) {
	sh := mocks.NewShellMock(t)
	ctx := internal.NewInternalContext(
		context.TODO(),
		mocks.NewFileSystemMock(t),
		sh,
		mocks.NewRepositoryMock(t),
		[]string{},
		ioutil.Discard,
	)

	actualSh := ctx.Shell()

	assert.Equal(t, sh, actualSh)
}

func TestContext_Repository(t *testing.T) {
	expected := mocks.NewRepositoryMock(t)
	ctx := internal.NewInternalContext(
		context.TODO(),
		mocks.NewFileSystemMock(t),
		mocks.NewShellMock(t),
		expected,
		[]string{},
		ioutil.Discard,
	)

	actual := ctx.Repository()

	assert.Equal(t, expected, actual)
}

func TestContext_Args(t *testing.T) {
	expected := []string{"param"}
	ctx := internal.NewInternalContext(
		context.TODO(),
		mocks.NewFileSystemMock(t),
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		expected,
		ioutil.Discard,
	)

	actual := ctx.Args()

	assert.Equal(t, expected, actual)
}

func TestContext_Output(t *testing.T) {
	expected := bytes.NewBufferString("")
	ctx := internal.NewInternalContext(
		context.TODO(),
		mocks.NewFileSystemMock(t),
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{},
		expected,
	)

	actual := ctx.Output()

	assert.Equal(t, expected, actual)
}

func TestContext_Message(t *testing.T) {
	expectedMessage := "MESSAGE"
	fs := mocks.NewFileSystemMock(t).
		ReadMock.When("filepath").Then(expectedMessage, nil)

	ctx := internal.NewInternalContext(
		context.TODO(),
		fs,
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{"filepath"},
		ioutil.Discard,
	)

	actual := ctx.Message()

	assert.Equal(t, expectedMessage, actual)
}

func TestContext_Message_Multiple(t *testing.T) {
	expectedMessage := "MESSAGE"
	fs := mocks.NewFileSystemMock(t).
		ReadMock.Expect("filepath").Return(expectedMessage, nil)

	ctx := internal.NewInternalContext(
		context.TODO(),
		fs,
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{"filepath"},
		ioutil.Discard,
	)

	for i := 0; i < 3; i++ {
		actual := ctx.Message()

		assert.Equal(t, expectedMessage, actual)
	}

	assert.Equal(t, 1, len(fs.ReadMock.Calls()))
}

func TestContext_Message_ErrorReading(t *testing.T) {
	ctx := internal.NewInternalContext(
		context.TODO(),
		mocks.NewFileSystemMock(t).
			ReadMock.Expect("filepath").Return("", errors.New("test")),
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{"filepath"},
		ioutil.Discard,
	)

	assert.Panics(t, func() {
		ctx.Message()
	})
}

func TestContext_Stop(t *testing.T) {
	ctx := internal.NewInternalContext(
		context.Background(),
		mocks.NewFileSystemMock(t),
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{},
		ioutil.Discard,
	)

	ctx.Stop()

	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestContext_Value(t *testing.T) {
	key := "this-is-key"
	expected := "this-is-value"

	ctx := internal.NewInternalContext(
		context.WithValue(context.Background(), key, expected), // nolint
		mocks.NewFileSystemMock(t),
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{},
		ioutil.Discard,
	)

	data := ctx.Value(key)

	assert.Equal(t, expected, data)
}

func TestContext_Deadline(t *testing.T) {
	ctx := internal.NewInternalContext(
		context.Background(),
		mocks.NewFileSystemMock(t),
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{},
		ioutil.Discard,
	)

	data, ok := ctx.Deadline()

	assert.NotNil(t, data)
	assert.False(t, ok)
}

func TestContext_Done(t *testing.T) {
	ctx := internal.NewInternalContext(
		context.Background(),
		mocks.NewFileSystemMock(t),
		mocks.NewShellMock(t),
		mocks.NewRepositoryMock(t),
		[]string{},
		ioutil.Discard,
	)

	chanell := ctx.Done()

	assert.NotNil(t, chanell)
}
