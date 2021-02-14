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

func TestAddToIndex_NotConfigured(t *testing.T) {
	rule := rules.AddToIndex{}

	err := rule.Check(mocks.NewExecutionContextMock(t), ioutil.Discard)

	assert.NoError(t, err)
}

func TestAddToIndex_CorrectAddToIndex(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		AddGlobMock.When("glob1/*.go").Then(nil).
		AddGlobMock.When("*.css").Then(nil).
		AddGlobMock.When("mocks").Then(nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	rule := rules.AddToIndex{
		Globs: []rules.Glob{
			{"glob1/*.go", true},
			{"*.css", true},
			{"mocks", true},
		},
	}

	err := rule.Check(ctx, ioutil.Discard)

	assert.NoError(t, err)
}

func TestAddToIndex_FailedAddToIndex(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		AddGlobMock.When("glob1/*.go").Then(nil).
		AddGlobMock.When("*.css").Then(errors.New("testError")).
		AddGlobMock.When("mocks").Then(nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	rule := rules.AddToIndex{
		Globs: []rules.Glob{
			{"glob1/*.go", true},
			{"*.css", true},
			{"mocks", true},
		},
	}
	err := rule.Check(ctx, ioutil.Discard)

	assert.Error(t, err, "testError")
}

func TestAddToIndex_FailedAddToIndexOptional(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		AddGlobMock.When("glob1/*.go").Then(nil).
		AddGlobMock.When("*.css").Then(git.ErrGlobNoMatches).
		AddGlobMock.When("mocks").Then(nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	tests := []struct {
		name       string
		isRequired bool
	}{
		{name: "Optional true", isRequired: false},
		{name: "Optional false", isRequired: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := rules.AddToIndex{
				Globs: []rules.Glob{
					{"glob1/*.go", tt.isRequired},
					{"*.css", tt.isRequired},
					{"mocks", tt.isRequired},
				},
			}

			err := rule.Check(ctx, ioutil.Discard)

			if !tt.isRequired {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err, git.ErrGlobNoMatches)
			}
		})
	}
}

func TestAddToIndex_GetPosition(t *testing.T) {
	rule := rules.AddToIndex{}

	assert.Equal(t, rules.PostScripts, rule.GetPosition())
}

func TestAddToIndex_Compile(t *testing.T) {
	rule := rules.AddToIndex{
		Globs: []rules.Glob{
			{Glob: "{{var1}}", IsRequired: false},
			{Glob: "data", IsRequired: false},
		},
	}

	rule.Compile(map[string]interface{}{"var1": "VALUE"})

	assert.Equal(t, rules.AddToIndex{
		Globs: []rules.Glob{
			{Glob: "VALUE", IsRequired: false},
			{Glob: "data", IsRequired: false},
		},
	}, rule)
}
