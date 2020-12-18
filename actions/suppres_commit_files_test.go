package actions_test

import (
	"errors"
	"fisherman/actions"
	"fisherman/mocks"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestSuppresCommitFiles_NotConfigured(t *testing.T) {
	next, err := actions.SuppresCommitFiles(mocks.NewSyncContextMock(t), actions.SuppresCommitFilesSections{})

	assert.NoError(t, err)
	assert.True(t, next)
}

func TestSuppresCommitFiles_WithRemoveFromIndex(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		RemoveGlobMock.When("glob1/demo.go").Then(nil).
		RemoveGlobMock.When("demo.css").Then(git.ErrGlobNoMatches).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go", "demo.css"}, nil)

	ctx := mocks.NewSyncContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.SuppresCommitFiles(ctx, actions.SuppresCommitFilesSections{
		Globs:           []string{"glob1/*.go", "*.css", "mocks"},
		RemoveFromIndex: true,
	})

	assert.NoError(t, err)
	assert.True(t, next)
}

func TestSuppresCommitFiles(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		RemoveGlobMock.When("glob1/demo.go").Then(nil).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

	ctx := mocks.NewSyncContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.SuppresCommitFiles(ctx, actions.SuppresCommitFilesSections{
		Globs: []string{"glob1/*.go", "*.css", "mocks"},
	})

	assert.EqualError(t, err, "1 error occurred:\n\t* file glob1/demo.go can not be committed\n\n")
	assert.True(t, next)
}

func TestSuppresCommitFiles_GetFilesInIndexError(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, errors.New("test error"))

	ctx := mocks.NewSyncContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.SuppresCommitFiles(ctx, actions.SuppresCommitFilesSections{
		Globs: []string{"glob1/*.go", "*.css", "mocks"},
	})

	assert.EqualError(t, err, "test error")
	assert.False(t, next)
}

func TestSuppresCommitFiles_MatchingError(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

	ctx := mocks.NewSyncContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.SuppresCommitFiles(ctx, actions.SuppresCommitFilesSections{
		Globs: []string{"[/"},
	})

	assert.EqualError(t, err, "syntax error in pattern")
	assert.False(t, next)
}

func TestSuppresCommitFiles_RemoveGlobError(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		RemoveGlobMock.Expect("glob1/demo.go").Return(errors.New("test error")).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

	ctx := mocks.NewSyncContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.SuppresCommitFiles(ctx, actions.SuppresCommitFilesSections{
		Globs:           []string{"glob1/*.go", "*.css", "mocks"},
		RemoveFromIndex: true,
	})

	assert.EqualError(t, err, "test error")
	assert.False(t, next)
}

func TestSuppresCommitFilesSections_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected *actions.SuppresCommitFilesSections
	}{
		{
			name:   "from strings",
			source: "[glob1, glob2, glob3]",
			expected: &actions.SuppresCommitFilesSections{
				Globs: []string{"glob1", "glob2", "glob3"},
			},
		},
		{
			name: "from structure",
			source: `
globs: [glob1, glob2, glob3]
remove-from-index: true
`,
			expected: &actions.SuppresCommitFilesSections{
				Globs:           []string{"glob1", "glob2", "glob3"},
				RemoveFromIndex: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual actions.SuppresCommitFilesSections
			err := yaml.Unmarshal([]byte(tt.source), &actual)

			assert.ObjectsAreEqual(tt.expected, actual)
			assert.NoError(t, err)
		})
	}
}
