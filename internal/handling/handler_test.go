package handling_test

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/expression"
	. "fisherman/internal/handling"
	"fisherman/internal/validation"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Add test for ctx Cancel

func TestHookHandler_Handle(t *testing.T) {
	ctx := mocks.NewExecutionContextMock(t).
		OutputMock.Return(testutils.NopCloser(io.Discard)).
		GlobalVariablesMock.Return(map[string]interface{}{}, nil)

	engine := expression.NewGoExpressionEngine()

	testCases := []struct {
		name         string
		engine       expression.Engine
		rules        []configuration.Rule
		workersCount int
		expectedErr  string
	}{
		{
			name:   "rules checked successfully",
			engine: engine,
			rules: []configuration.Rule{
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST1").
					GetContitionMock.Return("").
					GetTypeMock.Return("rule1").
					CheckMock.Return(nil),
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST2").
					GetContitionMock.Return("1 == 1").
					GetTypeMock.Return("rule2").
					CheckMock.Return(nil),
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST3").
					GetContitionMock.Return("1 == 3").
					GetTypeMock.Return("rule3").
					CheckMock.Return(nil),
			},
			workersCount: 2,
		},
		{
			name: "expression engin returns error",
			engine: mocks.NewEngineMock(t).
				EvalMock.Return(false, validation.Errorf("test-rule", "test")),
			rules: []configuration.Rule{
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST").
					GetContitionMock.Return("").
					GetTypeMock.Return("rule1").
					CheckMock.Return(nil),
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST").
					GetContitionMock.Return("1 == 1").
					GetTypeMock.Return("rule2"),
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST").
					GetContitionMock.Return("1 == 2").
					GetTypeMock.Return("rule3"),
			},
			workersCount: 2,
			expectedErr:  "[test-rule] test",
		},
		{
			name:   "rule returns error",
			engine: engine,
			rules: []configuration.Rule{
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST").
					GetContitionMock.Return("").
					GetTypeMock.Return("rule1").
					CheckMock.Return(nil),
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST").
					GetContitionMock.Return("1 == 1").
					GetTypeMock.Return("rule2").
					CheckMock.Return(validation.Errorf("test-rule", "test")),
				mocks.NewRuleMock(t).
					GetPrefixMock.Return("TEST").
					GetContitionMock.Return("1 == 3").
					GetTypeMock.Return("rule3").
					CheckMock.Return(nil),
			},
			workersCount: 2,
			expectedErr:  "1 error occurred:\n\t* [rule2] [test-rule] test\n\n",
		},
		{
			name:   "configuration error",
			engine: engine,
			rules: []configuration.Rule{
				mocks.NewRuleMock(t).GetContitionMock.Return(""),
				mocks.NewRuleMock(t).GetContitionMock.Return(""),
				mocks.NewRuleMock(t).GetContitionMock.Return(""),
			},
			workersCount: 0,
			expectedErr:  "incorrect workers count",
		},
	}

	t.Run("runs rules", func(t *testing.T) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				handler := HookHandler{
					Engine:       tt.engine,
					Rules:        tt.rules,
					WorkersCount: tt.workersCount,
				}

				err := handler.Handle(ctx, []string{})

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("runs scripts", func(t *testing.T) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				handler := HookHandler{
					Engine:       tt.engine,
					Scripts:      tt.rules,
					WorkersCount: tt.workersCount,
				}

				err := handler.Handle(ctx, []string{})

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("runs post scripts", func(t *testing.T) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				handler := HookHandler{
					Engine:          tt.engine,
					PostScriptRules: tt.rules,
					WorkersCount:    tt.workersCount,
				}

				err := handler.Handle(ctx, []string{})

				testutils.AssertError(t, tt.expectedErr, err)
				print(tt.name)
			})
		}
	})

	t.Run("run rules in correct orders", func(t *testing.T) {
		order := []string{}

		checkHandler := func(ruleType string) func(internal.ExecutionContext, io.Writer) error {
			return func(ec internal.ExecutionContext, w io.Writer) error {
				order = append(order, ruleType)

				return nil
			}
		}

		handler := HookHandler{
			Engine: mocks.NewEngineMock(t),
			Rules: []configuration.Rule{
				mocks.NewRuleMock(t).
					GetContitionMock.Return("").
					GetPrefixMock.Return("rule").
					CheckMock.Set(checkHandler("rule")),
			},
			Scripts: []configuration.Rule{
				mocks.NewRuleMock(t).
					GetContitionMock.Return("").
					GetPrefixMock.Return("script").
					CheckMock.Set(checkHandler("script")),
			},
			PostScriptRules: []configuration.Rule{
				mocks.NewRuleMock(t).
					GetContitionMock.Return("").
					GetPrefixMock.Return("post-script").
					CheckMock.Set(checkHandler("post-script")),
			},
			WorkersCount: 10,
		}

		err := handler.Handle(ctx, []string{})

		assert.NoError(t, err)
		assert.Equal(t, []string{"rule", "script", "post-script"}, order)
	})
}
