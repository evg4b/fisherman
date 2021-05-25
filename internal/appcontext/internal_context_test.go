package appcontext_test

import (
	"bytes"
	"context"
	"errors"
	infra "fisherman/internal"
	"fisherman/internal/appcontext"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_Files(t *testing.T) {
	expectedFs := mocks.NewFileSystemMock(t)
	ctx := appcontext.NewContextBuilder().
		WithFileSystem(expectedFs).
		WithShell(mocks.NewShellMock(t)).
		WithRepository(mocks.NewRepositoryMock(t)).
		Build()

	actualFs := ctx.Files()

	assert.Equal(t, expectedFs, actualFs)
}

func TestContext_Shell(t *testing.T) {
	expectedShell := mocks.NewShellMock(t)
	ctx := appcontext.NewContextBuilder().
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithShell(expectedShell).
		WithRepository(mocks.NewRepositoryMock(t)).
		Build()

	actualSh := ctx.Shell()

	assert.Equal(t, expectedShell, actualSh)
}

func TestContext_Repository(t *testing.T) {
	expectedRepo := mocks.NewRepositoryMock(t)
	ctx := appcontext.NewContextBuilder().
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithShell(mocks.NewShellMock(t)).
		WithRepository(expectedRepo).
		Build()

	actualRepo := ctx.Repository()

	assert.Equal(t, expectedRepo, actualRepo)
}

func TestContext_Args(t *testing.T) {
	expectedArgs := []string{"param"}
	ctx := appcontext.NewContextBuilder().
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithShell(mocks.NewShellMock(t)).
		WithRepository(mocks.NewRepositoryMock(t)).
		WithArgs(expectedArgs).
		Build()

	actualArgs := ctx.Args()

	assert.Equal(t, expectedArgs, actualArgs)
}

func TestContext_Output(t *testing.T) {
	expectedOutput := bytes.NewBufferString("")
	ctx := appcontext.NewContextBuilder().
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithShell(mocks.NewShellMock(t)).
		WithRepository(mocks.NewRepositoryMock(t)).
		WithOutput(expectedOutput).
		Build()

	actualOutput := ctx.Output()

	assert.Equal(t, expectedOutput, actualOutput)
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
			ctx := appcontext.NewContextBuilder().
				WithFileSystem(testutils.FsFromMap(t, tt.files)).
				WithShell(mocks.NewShellMock(t)).
				WithRepository(mocks.NewRepositoryMock(t)).
				WithArgs(tt.args).
				Build()

			actual, err := ctx.Message()

			assert.Equal(t, tt.expected, actual)
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}

func TestContext_Stop(t *testing.T) {
	ctx := appcontext.NewContextBuilder().
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithShell(mocks.NewShellMock(t)).
		WithRepository(mocks.NewRepositoryMock(t)).
		Build()

	ctx.Stop()

	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestContext_Value(t *testing.T) {
	key := "this-is-key"
	expected := "this-is-value"

	ctx := appcontext.NewContextBuilder().
		WithContext(context.WithValue(context.Background(), key, expected)). //nolint
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithShell(mocks.NewShellMock(t)).
		WithRepository(mocks.NewRepositoryMock(t)).
		Build()

	data := ctx.Value(key)

	assert.Equal(t, expected, data)
}

func TestContext_Deadline(t *testing.T) {
	ctx := appcontext.NewContextBuilder().
		WithContext(context.Background()).
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithShell(mocks.NewShellMock(t)).
		WithRepository(mocks.NewRepositoryMock(t)).
		Build()

	data, ok := ctx.Deadline()

	assert.NotNil(t, data)
	assert.False(t, ok)
}

func TestContext_Done(t *testing.T) {
	ctx := appcontext.NewContextBuilder().
		WithContext(context.Background()).
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithShell(mocks.NewShellMock(t)).
		WithRepository(mocks.NewRepositoryMock(t)).
		Build()

	chanell := ctx.Done()

	assert.NotNil(t, chanell)
}

func TestContext_GlobalVariables(t *testing.T) {
	tests := []struct {
		name        string
		expected    map[string]interface{}
		repository  infra.Repository
		expectedErr string
	}{
		{
			name: "GetLastTag returns error",
			repository: mocks.NewRepositoryMock(t).
				GetLastTagMock.Return("", errors.New("GetLastTag error")),
			expected:    nil,
			expectedErr: "GetLastTag error",
		},
		{
			name: "GetCurrentBranch returns error",
			repository: mocks.NewRepositoryMock(t).
				GetLastTagMock.Return("1.0.0", nil).
				GetCurrentBranchMock.Return("", errors.New("GetCurrentBranch error")),
			expected:    nil,
			expectedErr: "GetCurrentBranch error",
		},
		{
			name: "GetUser returns error",
			repository: mocks.NewRepositoryMock(t).
				GetLastTagMock.Return("1.0.0", nil).
				GetCurrentBranchMock.Return("refs/head/develop", nil).
				GetUserMock.Return(infra.User{}, errors.New("GetUser error")),
			expected:    nil,
			expectedErr: "GetUser error",
		},
		{
			name: "GetUser returns error",
			repository: mocks.NewRepositoryMock(t).
				GetLastTagMock.Return("1.0.0", nil).
				GetCurrentBranchMock.Return("refs/head/develop", nil).
				GetUserMock.Return(infra.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil),
			expected: map[string]interface{}{
				"Tag":        "1.0.0",
				"BranchName": "refs/head/develop",
				"UserName":   "evg4b",
				"UserEmail":  "evg4b@mail.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := appcontext.NewContextBuilder().
				WithFileSystem(mocks.NewFileSystemMock(t)).
				WithShell(mocks.NewShellMock(t)).
				WithRepository(tt.repository).
				Build()

			actual, err := ctx.GlobalVariables()

			assert.Equal(t, tt.expected, actual)
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}
