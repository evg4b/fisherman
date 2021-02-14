package rules_test

import (
	"errors"
	"fisherman/internal/rules"
	"fisherman/testing/mocks"
	"io/ioutil"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestSuppressCommitFiles_GetPosition(t *testing.T) {
	rule := rules.SuppressCommitFiles{
		BaseRule:        rules.BaseRule{Type: "demo-rule", Condition: "rule-condition"},
		Globs:           []string{"glob1", "glob2", "glob3"},
		RemoveFromIndex: true,
	}

	assert.Equal(t, rules.PostScripts, rule.GetPosition())
}

func TestSuppresCommitFiles_NotConfigured(t *testing.T) {
	ctx := mocks.NewExecutionContextMock(t)
	rule := rules.SuppressCommitFiles{}

	err := rule.Check(ctx, ioutil.Discard)

	assert.NoError(t, err)
}

func TestSuppresCommitFiles(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		RemoveGlobMock.When("glob1/demo.go").Then(nil).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	rule := rules.SuppressCommitFiles{
		Globs: []string{"glob1/*.go", "*.css", "mocks"},
	}

	err := rule.Check(ctx, ioutil.Discard)

	assert.EqualError(t, err, "1 error occurred:\n\t* file glob1/demo.go can not be committed\n\n")
}

func TestSuppresCommitFiles_WithRemoveFromIndex(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		RemoveGlobMock.When("glob1/demo.go").Then(nil).
		RemoveGlobMock.When("demo.css").Then(git.ErrGlobNoMatches).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go", "demo.css"}, nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	rule := rules.SuppressCommitFiles{
		Globs:           []string{"glob1/*.go", "*.css", "mocks"},
		RemoveFromIndex: true,
	}

	err := rule.Check(ctx, ioutil.Discard)

	assert.NoError(t, err)
}

func TestSuppresCommitFiles_GetFilesInIndexError(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, errors.New("test error"))

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	rule := rules.SuppressCommitFiles{
		Globs: []string{"glob1/*.go", "*.css", "mocks"},
	}

	err := rule.Check(ctx, ioutil.Discard)

	assert.EqualError(t, err, "test error")
}

func TestSuppresCommitFiles_MatchingError(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	rule := rules.SuppressCommitFiles{
		Globs: []string{"[/"},
	}

	err := rule.Check(ctx, ioutil.Discard)

	assert.EqualError(t, err, "syntax error in pattern")
}

func TestSuppresCommitFiles_RemoveGlobError(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		RemoveGlobMock.Expect("glob1/demo.go").Return(errors.New("test error")).
		GetFilesInIndexMock.Expect().Return([]string{"glob1/demo.go"}, nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	rule := rules.SuppressCommitFiles{
		Globs:           []string{"glob1/*.go", "*.css", "mocks"},
		RemoveFromIndex: true,
	}

	err := rule.Check(ctx, ioutil.Discard)

	assert.EqualError(t, err, "test error")
}

func TestSuppressCommitFiles_Compile(t *testing.T) {
	rule := rules.SuppressCommitFiles{
		Globs:           []string{"{{var1}}", "{{var1}}.css", "mocks"},
		RemoveFromIndex: true,
		BaseRule: rules.BaseRule{
			Type:      "{{var1}}",
			Condition: "{{var1}}",
		},
	}

	rule.Compile(map[string]interface{}{"var1": "VALUE"})

	assert.Equal(t, rules.SuppressCommitFiles{
		Globs:           []string{"VALUE", "VALUE.css", "mocks"},
		RemoveFromIndex: true,
		BaseRule: rules.BaseRule{
			Type:      "{{var1}}",
			Condition: "VALUE",
		},
	}, rule)
}
