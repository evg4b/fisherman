package configuration_test

import (
	. "fisherman/configuration"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariablesSection_Compile_Empty(t *testing.T) {
	section := VariablesSection{}

	assert.NotPanics(t, func() {
		variables, err := section.Compile(mocks.NewEngineMock(t), map[string]interface{}{})

		assert.Empty(t, variables)
		assert.NoError(t, err)
	})
}

func TestVariables_Compile(t *testing.T) {
	engine := mocks.NewEngineMock(t).EvalMapMock.Return(map[string]interface{}{}, nil)

	section := VariablesSection{
		StaticVariables: map[string]string{
			"VAR_1": "{{var1}}",
			"VAR_2": "{{var2}}_demo",
			"VAR_3": "{var2}_test",
		},
		ExtractVariables: []string{
			"Extract({{var1}}, {{var2}})",
			"Extract('{{var1}}', \"{{var1}}\")",
		},
	}

	_, err := section.Compile(engine, map[string]interface{}{
		"var1": "localValue1",
		"var2": "localValue2",
	})

	assert.NoError(t, err)

	assert.Equal(t, map[string]string{
		"VAR_1": "localValue1",
		"VAR_2": "localValue2_demo",
		"VAR_3": "{var2}_test",
	}, section.StaticVariables)

	assert.Equal(t, []string{
		"Extract(localValue1, localValue2)",
		"Extract('localValue1', \"localValue1\")",
	}, section.ExtractVariables)
}

func TestVariablesSection_CompileAndReturnVariables(t *testing.T) {
	section := VariablesSection{
		ExtractVariables: []string{"stub"},
	}
	engine := mocks.NewEngineMock(t).EvalMapMock.Return(map[string]interface{}{
		"var1": "new value",
	}, nil)

	assert.NotPanics(t, func() {
		variables, err := section.Compile(engine, map[string]interface{}{
			"var1": "value",
			"var2": "value2",
		})

		assert.Equal(t, map[string]interface{}{
			"var1": "new value",
			"var2": "value2",
		}, variables)
		assert.NoError(t, err)
	})
}
