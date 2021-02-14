package configuration_test

import (
	"errors"
	"fisherman/configuration"
	"fisherman/internal/expression"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"testing"
)

func TestCommonConfig_Compile(t *testing.T) {
	baseEngine := mocks.NewEngineMock(t).EvalMapMock.Return(map[string]interface{}{
		"VAR1": "THIS IS 2 VAR",
	}, nil)

	tests := []struct {
		name          string
		config        *configuration.CommonConfig
		expectedError string
		engine        expression.Engine
	}{
		{
			name:   "empty rule",
			config: &configuration.CommonConfig{},
			engine: baseEngine,
		},
		{
			name: "",
			config: &configuration.CommonConfig{
				VariablesSection: configuration.VariablesSection{
					StaticVariables: map[string]string{
						"VAR1": "%{{VAR1}}%",
					},
					ExtractVariables: []string{
						"Stub({{VAR1}})",
					},
				},
				RulesSection: configuration.RulesSection{
					Rules: []configuration.Rule{
						mocks.NewRuleMock(t).CompileMock.Return(),
					},
				},
			},
			engine:        baseEngine,
			expectedError: "",
		},
		{
			name: "",
			config: &configuration.CommonConfig{
				VariablesSection: configuration.VariablesSection{
					StaticVariables: map[string]string{
						"VAR1": "%{{VAR1}}%",
					},
					ExtractVariables: []string{
						"Stub({{VAR1}})",
					},
				},
				RulesSection: configuration.RulesSection{
					Rules: []configuration.Rule{
						mocks.NewRuleMock(t).CompileMock.Return(),
					},
				},
			},
			engine:        mocks.NewEngineMock(t).EvalMapMock.Return(map[string]interface{}{}, errors.New("test error")),
			expectedError: "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Compile(tt.engine, map[string]interface{}{})

			testutils.CheckError(t, tt.expectedError, err)
		})
	}
}
