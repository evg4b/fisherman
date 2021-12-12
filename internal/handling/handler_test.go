package handling_test

import (
	"context"
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
	"fisherman/internal/expression"
	. "fisherman/internal/handling"
	"fisherman/internal/rules"
	"fisherman/internal/validation"
	"fisherman/pkg/guards"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHookHandler_Handle(t *testing.T) {
	engine := expression.NewGoExpressionEngine()

	validRule := mocks.NewRuleMock(t).
		GetTypeMock.Return(rules.ExecType).
		InitMock.Return().
		GetPositionMock.Return(rules.Scripts).
		CompileMock.Return().
		GetContitionMock.Return("").
		GetPrefixMock.Return("prefix-").
		CheckMock.Return(nil)

	testCases := []struct {
		name         string
		engine       expression.Engine
		rules        []*mocks.RuleMock
		workersCount uint
		expectedErr  string
	}{
		{
			name:   "rules checked successfully",
			engine: engine,
			rules: []*mocks.RuleMock{
				ruleMock(t, "TEST1", "", rules.ShellScriptType).
					InitMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST2", "1 == 1", rules.ExecType).
					InitMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST3", "1 == 3", rules.SuppressCommitFilesType).
					InitMock.Return().
					CheckMock.Return(nil),
			},
			workersCount: 2,
		},
		{
			name: "expression engin returns error",
			engine: mocks.NewEngineMock(t).
				EvalMock.Return(false, validation.Errorf("test-rule", "test")),
			rules: []*mocks.RuleMock{
				ruleMock(t, "TEST1", "", rules.ShellScriptType).
					InitMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST", "1 == 1", rules.ExecType).
					InitMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST", "1 == 2", rules.SuppressCommitFilesType).
					InitMock.Return().
					CheckMock.Return(nil),
			},
			workersCount: 2,
			expectedErr:  "[test-rule] test",
		},
		{
			name:   "rule returns error",
			engine: engine,
			rules: []*mocks.RuleMock{
				ruleMock(t, "TEST", "", rules.ShellScriptType).
					InitMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST", "1 == 1", rules.ExecType).
					InitMock.Return().
					CheckMock.Return(validation.Errorf("test-rule", "test")),
				ruleMock(t, "TEST", "1 == 3", rules.SuppressCommitFilesType).
					InitMock.Return().
					CheckMock.Return(nil),
			},
			workersCount: 2,
			expectedErr:  "1 error occurred:\n\t* [exec] [test-rule] test\n\n",
		},
		{
			name:   "configuration error",
			engine: engine,
			rules: []*mocks.RuleMock{
				ruleMock(t, "TEST", "", rules.ShellScriptType),
				ruleMock(t, "TEST", "", rules.ShellScriptType),
				ruleMock(t, "TEST", "", rules.ShellScriptType),
			},
			workersCount: 0,
			expectedErr:  "incorrect workers count",
		},
	}

	t.Run("runs rules", func(t *testing.T) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				handler, err := NewHookHandler(
					"pre-commit",
					WithExpressionEngine(tt.engine),
					WithHooksConfig(&configuration.HooksConfig{
						PreCommitHook: &configuration.HookConfig{
							Rules: applyPosition(tt.rules, rules.PreScripts),
						},
					}),
					WithWorkersCount(tt.workersCount),
				)

				guards.NoError(err)

				err = handler.Handle(context.TODO())

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("runs scripts", func(t *testing.T) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				handler, err := NewHookHandler(
					"pre-commit",
					WithExpressionEngine(tt.engine),
					WithHooksConfig(&configuration.HooksConfig{
						PreCommitHook: &configuration.HookConfig{
							Rules: applyPosition(tt.rules, rules.Scripts),
						},
					}),
					WithWorkersCount(tt.workersCount),
				)

				guards.NoError(err)

				err = handler.Handle(context.TODO())

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("runs post scripts", func(t *testing.T) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				handler, err := NewHookHandler(
					"pre-commit",
					WithExpressionEngine(tt.engine),
					WithHooksConfig(&configuration.HooksConfig{
						PreCommitHook: &configuration.HookConfig{
							Rules: applyPosition(tt.rules, rules.PostScripts),
						},
					}),
					WithWorkersCount(tt.workersCount),
				)

				guards.NoError(err)

				err = handler.Handle(context.TODO())

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("run rules in correct orders", func(t *testing.T) {
		order := []string{}

		ruleDef := func(ruleType string, position byte) configuration.Rule {
			checkHandler := func() func(context.Context, io.Writer) error {
				return func(ec context.Context, w io.Writer) error {
					order = append(order, ruleType)

					return nil
				}
			}

			return mocks.NewRuleMock(t).
				GetTypeMock.Return(rules.ExecType).
				InitMock.Return().
				CompileMock.Return().
				GetContitionMock.Return("").
				GetPrefixMock.Return("prefix-").
				CheckMock.Set(checkHandler()).
				GetPositionMock.Return(position)
		}

		handler, err := NewHookHandler(
			"pre-commit",
			WithExpressionEngine(mocks.NewEngineMock(t)),
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: []configuration.Rule{
						ruleDef("rule", rules.PreScripts),
						ruleDef("script", rules.Scripts),
						ruleDef("post-script", rules.PostScripts),
					},
				},
			}),
		)

		guards.NoError(err)

		err = handler.Handle(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, []string{"rule", "script", "post-script"}, order)
	})

	t.Run("canceled context", func(t *testing.T) {
		handler, err := NewHookHandler(
			"pre-commit",
			WithExpressionEngine(mocks.NewEngineMock(t)),
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: []configuration.Rule{
						validRule,
					},
				},
			}),
			WithWorkersCount(4),
			WithFileSystem(mocks.NewFilesystemMock(t)),
			WithRepository(mocks.NewRepositoryMock(t)),
		)

		guards.NoError(err)

		ctx, cancel := context.WithCancel(context.Background())

		cancel()

		err = handler.Handle(ctx)

		assert.EqualError(t, err, "1 error occurred:\n\t* context canceled\n\n")
	})

	t.Run("hook config is not presented", func(t *testing.T) {
		_, err := NewHookHandler(
			constants.PreCommitHook,
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: nil,
			}),
		)

		assert.Error(t, err, ErrNotPresented.Error())
	})

	t.Run("invalid hook name", func(t *testing.T) {
		_, err := NewHookHandler(
			"unknown-hook",
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: nil,
			}),
		)

		assert.Error(t, err, "'unknown-hook' is not valid hook name")
	})
}

func ruleMock(t *testing.T, prefix, condition, ruleType string) *mocks.RuleMock {
	t.Helper()

	return mocks.NewRuleMock(t).
		GetPrefixMock.Return(prefix).
		GetContitionMock.Return(condition).
		GetTypeMock.Return(ruleType).
		CompileMock.Return().
		InitMock.Return()
}

func applyPosition(rules []*mocks.RuleMock, position byte) []configuration.Rule {
	updatedRules := []configuration.Rule{}
	for _, rule := range rules {
		rule.GetPositionMock.Return(position)
		updatedRules = append(updatedRules, rule)
	}

	return updatedRules
}
