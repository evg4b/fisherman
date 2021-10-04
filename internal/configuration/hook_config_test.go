package configuration_test

import (
	. "fisherman/internal/configuration"
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHookConfig_Compile_Empty(t *testing.T) {
	section := HookConfig{}

	assert.NotPanics(t, func() {
		variables, err := section.Compile(map[string]interface{}{})

		assert.Empty(t, variables)
		assert.NoError(t, err)
	})
}

func TestHookConfig_VariablesSections_Compile(t *testing.T) {
	section := HookConfig{
		StaticVariables: map[string]string{
			"VAR_1": "{{var1}}",
			"VAR_2": "{{var2}}_demo",
			"VAR_3": "{var2}_test",
		},
	}

	_, err := section.Compile(map[string]interface{}{
		"var1": "localValue1",
		"var2": "localValue2",
	})

	assert.NoError(t, err)

	assert.Equal(t, map[string]string{
		"VAR_1": "localValue1",
		"VAR_2": "localValue2_demo",
		"VAR_3": "{var2}_test",
	}, section.StaticVariables)
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
  - source: demo
    expression: expr
`,
			expected: HookConfig{
				StaticVariables: map[string]string{
					"demo":  "Test",
					"demo2": "Test2",
				},
				ExtractVariables: []ExtractVariable{
					{"demo", "expr"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var section HookConfig

			err := testutils.DecodeYaml(tt.source, &section)

			assert.ObjectsAreEqual(tt.expected, section)
			testutils.CheckError(t, tt.expectedError, err)
		})
	}
}
