package configuration_test

import (
	"errors"
	. "fisherman/configuration"
	"fisherman/internal/expression"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestHookConfig_Compile(t *testing.T) {
	baseEngine := mocks.NewEngineMock(t).EvalMapMock.Return(map[string]interface{}{
		"VAR1": "THIS IS 2 VAR",
	}, nil)

	tests := []struct {
		name          string
		config        *HookConfig
		expectedError string
		engine        expression.Engine
	}{
		{
			name:   "empty rule",
			config: &HookConfig{},
			engine: baseEngine,
		},
		{
			name: "",
			config: &HookConfig{
				StaticVariables: map[string]string{
					"VAR1": "%{{VAR1}}%",
				},
				ExtractVariables: []string{
					"Stub({{VAR1}})",
				},
				Rules: []Rule{
					mocks.NewRuleMock(t).CompileMock.Return(),
				},
			},
			engine:        baseEngine,
			expectedError: "",
		},
		{
			name: "",
			config: &HookConfig{
				StaticVariables: map[string]string{
					"VAR1": "%{{VAR1}}%",
				},
				ExtractVariables: []string{
					"Stub({{VAR1}})",
				},
				Rules: []Rule{
					mocks.NewRuleMock(t).CompileMock.Return(),
				},
			},
			engine:        mocks.NewEngineMock(t).EvalMapMock.Return(map[string]interface{}{}, errors.New("test error")),
			expectedError: "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.config.Compile(tt.engine, map[string]interface{}{})

			testutils.CheckError(t, tt.expectedError, err)
		})
	}
}

func TestHookConfig_Compile_Empty(t *testing.T) {
	section := HookConfig{}

	assert.NotPanics(t, func() {
		variables, err := section.Compile(mocks.NewEngineMock(t), map[string]interface{}{})

		assert.Empty(t, variables)
		assert.NoError(t, err)
	})
}

func TestHookConfig_VariablesSections_Compile(t *testing.T) {
	engine := mocks.NewEngineMock(t).EvalMapMock.Return(map[string]interface{}{}, nil)

	section := HookConfig{
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

func TestHookConfig_CompileAndReturnVariables(t *testing.T) {
	section := HookConfig{
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

func TestHookConfig_VariablesSections_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      HookConfig
		expectedError string
	}{
		{
			name: "test1",
			source: `
variables:
  demo: Test
  demo2: Test2
extract-variables:
  - Extract("", "")
`,
			expected: HookConfig{
				StaticVariables: map[string]string{
					"demo":  "Test",
					"demo2": "Test2",
				},
				ExtractVariables: []string{
					"Extract(\"\", \"\")",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.source)
			decoder := yaml.NewDecoder(reader)
			decoder.KnownFields(true)

			var section HookConfig

			err := decoder.Decode(&section)

			assert.ObjectsAreEqual(tt.expected, section)
			testutils.CheckError(t, tt.expectedError, err)
		})
	}
}
