package rules_test

import (
	"context"
	"errors"
	"github.com/evg4b/fisherman/testing/mocks"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/evg4b/fisherman/internal/rules"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestSuppressCommitFiles_GetPosition(t *testing.T) {
	rule := SuppressCommitFiles{
		BaseRule:        BaseRule{Type: SuppressCommitFilesType},
		Globs:           []string{"glob1", "glob2", "glob3"},
		RemoveFromIndex: true,
	}

	assert.Equal(t, PostScripts, rule.GetPosition())
}

func TestSuppressCommitFiles_Compile(t *testing.T) {
	rule := SuppressCommitFiles{
		BaseRule: BaseRule{
			Type:      SuppressCommitFilesType,
			Condition: "{{var1}}",
		},
		Globs:           []string{"{{var1}}", "{{var1}}.css", "mocks"},
		RemoveFromIndex: true,
	}

	rule.Compile(map[string]any{"var1": "VALUE"})

	assert.Equal(t, SuppressCommitFiles{
		BaseRule: BaseRule{
			Type:      SuppressCommitFilesType,
			Condition: "VALUE",
		},
		Globs:           []string{"VALUE", "VALUE.css", "mocks"},
		RemoveFromIndex: true,
	}, rule)
}

func TestSuppressCommitFiles_Check(t *testing.T) {
	t.Run("not configured rule", func(t *testing.T) {
		rule := SuppressCommitFiles{
			BaseRule: BaseRule{Type: SuppressCommitFilesType},
		}

		err := rule.Check(context.TODO(), io.Discard)

		require.NoError(t, err)
	})

	t.Run("suppressed add files", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t).
			RemoveGlobMock.When("glob1/demo.go").Then(nil).
			GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

		rule := makeRule(
			&SuppressCommitFiles{
				BaseRule: BaseRule{Type: SuppressCommitFilesType},
				Globs:    []string{"glob1/*.go", "*.css", "mocks"},
			},
			WithRepository(repo),
		)

		err := rule.Check(context.TODO(), io.Discard)

		require.EqualError(t, err, "1 error occurred:\n\t* [suppress-commit-files] file glob1/demo.go can not be committed\n\n")
	})

	t.Run("removed files from index", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t).
			RemoveGlobMock.When("glob1/demo.go").Then(nil).
			RemoveGlobMock.When("demo.css").Then(git.ErrGlobNoMatches).
			GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go", "demo.css"}, nil)

		rule := makeRule(
			&SuppressCommitFiles{
				BaseRule:        BaseRule{Type: SuppressCommitFilesType},
				Globs:           []string{"glob1/*.go", "*.css", "mocks"},
				RemoveFromIndex: true,
			},
			WithRepository(repo),
		)

		err := rule.Check(context.TODO(), io.Discard)

		require.NoError(t, err)
	})

	t.Run("get files in index error", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t).
			GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, errors.New("test error"))

		rule := makeRule(
			&SuppressCommitFiles{
				BaseRule: BaseRule{Type: SuppressCommitFilesType},
				Globs:    []string{"glob1/*.go", "*.css", "mocks"},
			},
			WithRepository(repo),
		)

		err := rule.Check(context.TODO(), io.Discard)

		require.EqualError(t, err, "test error")
	})

	t.Run("glob patter parsing error", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t).
			GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

		rule := makeRule(
			&SuppressCommitFiles{
				BaseRule: BaseRule{Type: SuppressCommitFilesType},
				Globs:    []string{"[/"},
			},
			WithRepository(repo),
		)

		err := rule.Check(context.TODO(), io.Discard)

		require.EqualError(t, err, "syntax error in pattern")
	})

	t.Run("removing files from index error", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t).
			RemoveGlobMock.Expect("glob1/demo.go").Return(errors.New("test error")).
			GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

		rule := makeRule(
			&SuppressCommitFiles{
				BaseRule:        BaseRule{Type: SuppressCommitFilesType},
				Globs:           []string{"glob1/*.go", "*.css", "mocks"},
				RemoveFromIndex: true,
			},
			WithRepository(repo),
		)

		err := rule.Check(context.TODO(), io.Discard)

		require.EqualError(t, err, "test error")
	})
}
