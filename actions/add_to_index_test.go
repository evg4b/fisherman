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

func TestAddToIndex_NotConfigured(t *testing.T) {
	next, err := actions.AddToIndex(mocks.NewExecutionContextMock(t), []actions.Glob{})

	assert.NoError(t, err)
	assert.True(t, next)
}

func TestAddToIndex_CorrectAddToIndex(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		AddGlobMock.When("glob1/*.go").Then(nil).
		AddGlobMock.When("*.css").Then(nil).
		AddGlobMock.When("mocks").Then(nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.AddToIndex(ctx, []actions.Glob{
		{"glob1/*.go", true},
		{"*.css", true},
		{"mocks", true},
	})

	assert.NoError(t, err)
	assert.True(t, next)
}

func TestAddToIndex_FailedAddToIndex(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		AddGlobMock.When("glob1/*.go").Then(nil).
		AddGlobMock.When("*.css").Then(errors.New("testError")).
		AddGlobMock.When("mocks").Then(nil)

	ctx := mocks.NewExecutionContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.AddToIndex(ctx, []actions.Glob{
		{"glob1/*.go", true},
		{"*.css", true},
		{"mocks", true},
	})

	assert.Error(t, err, "testError")
	assert.False(t, next)
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
			next, err := actions.AddToIndex(ctx, []actions.Glob{
				{"glob1/*.go", tt.isRequired},
				{"*.css", tt.isRequired},
				{"mocks", tt.isRequired},
			})

			if !tt.isRequired {
				assert.NoError(t, err)
				assert.True(t, next)
			} else {
				assert.Equal(t, err, git.ErrGlobNoMatches)
				assert.False(t, next)
			}
		})
	}
}

func TestGlob_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name       string
		yamlSource string
		expected   actions.Glob
	}{
		{
			name:       "unmarshal string",
			yamlSource: "\"test string\"",
			expected:   actions.Glob{Glob: "test string", IsRequired: true},
		},
		{
			name: "unmarshal not required glob",
			yamlSource: `
glob: test glob structure
required: false
`,
			expected: actions.Glob{Glob: "test glob structure", IsRequired: false},
		},
		{
			name: "unmarshal not required glob",
			yamlSource: `
glob: another test glob structure
required: true
`,
			expected: actions.Glob{Glob: "another test glob structure", IsRequired: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result actions.Glob
			err := yaml.Unmarshal([]byte(tt.yamlSource), &result)

			assert.NoError(t, err)
			assert.ObjectsAreEqual(tt.expected, result)
		})
	}
}
