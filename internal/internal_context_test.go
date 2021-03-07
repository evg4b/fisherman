package internal_test

import (
	"bytes"
	"context"
	"fisherman/internal"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
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
	tests := []struct {
		name        string
		files       map[string]string
		expected    string
		expectedErr string
		args        []string
	}{
		{
			name:        "return message from file",
			files:       map[string]string{"filepath": "expectedMessage"},
			expected:    "expectedMessage",
			expectedErr: "",
			args:        []string{"filepath"},
		},
		{
			name:        "return message from file2",
			files:       map[string]string{},
			expected:    "",
			expectedErr: "argument at index 0 is not provided",
			args:        []string{},
		},
		{
			name:        "return message from file",
			files:       map[string]string{},
			expected:    "",
			expectedErr: "open filepath: file does not exist",
			args:        []string{"filepath"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := internal.NewInternalContext(
				context.TODO(),
				testutils.FsFromMap(t, tt.files),
				mocks.NewShellMock(t),
				mocks.NewRepositoryMock(t),
				tt.args,
				ioutil.Discard,
			)

			actual, err := ctx.Message()

			assert.Equal(t, tt.expected, actual)
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
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
