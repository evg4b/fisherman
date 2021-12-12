package handling_test

import (
	"context"
	"fisherman/internal/configuration"
	"fisherman/internal/expression"
	"fisherman/internal/handling"
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

	testCases := []struct {
		name         string
		engine       expression.Engine
		rules        []configuration.Rule
		workersCount uint
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

				handler, err := handling.NewHookHandler(
					"pre-commit",
					WithExpressionEngine(engine),
					WithHooksConfig(&configuration.HooksConfig{
						PreCommitHook: &configuration.HookConfig{
							Rules: tt.rules,
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
				handler, err := handling.NewHookHandler(
					"pre-commit",
					WithExpressionEngine(tt.engine),
					WithHooksConfig(&configuration.HooksConfig{
						PreCommitHook: &configuration.HookConfig{
							Rules: tt.rules,
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
				handler, err := handling.NewHookHandler(
					"pre-commit",
					WithExpressionEngine(tt.engine),
					WithHooksConfig(&configuration.HooksConfig{
						PreCommitHook: &configuration.HookConfig{
							Rules: tt.rules,
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

		checkHandler := func(ruleType string) func(context.Context, io.Writer) error {
			return func(ec context.Context, w io.Writer) error {
				order = append(order, ruleType)

				return nil
			}
		}

		handler, err := handling.NewHookHandler(
			"pre-commit",
			WithExpressionEngine(mocks.NewEngineMock(t)),
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: []configuration.Rule{
						mocks.NewRuleMock(t).
							GetContitionMock.Return("").
							GetPrefixMock.Return("rule").
							CheckMock.Set(checkHandler("rule")).
							GetPositionMock.Return(rules.PreScripts),
						mocks.NewRuleMock(t).
							GetContitionMock.Return("").
							GetPrefixMock.Return("script").
							CheckMock.Set(checkHandler("script")).
							GetPositionMock.Return(rules.Scripts),
						mocks.NewRuleMock(t).
							GetContitionMock.Return("").
							GetPrefixMock.Return("post-script").
							CheckMock.Set(checkHandler("post-script")).
							GetPositionMock.Return(rules.PostScripts),
					},
				},
			}),
			WithWorkersCount(10),
		)

		guards.NoError(err)

		err = handler.Handle(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, []string{"rule", "script", "post-script"}, order)
	})

	t.Run("canceled context", func(t *testing.T) {
		randomRule := func() configuration.Rule {
			return mocks.NewRuleMock(t).
				GetContitionMock.Return("").
				GetPrefixMock.Return("test-").
				GetPositionMock.Return(rules.Scripts)
		}

		handler, err := handling.NewHookHandler(
			"pre-commit",
			WithExpressionEngine(mocks.NewEngineMock(t)),
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: []configuration.Rule{
						randomRule(),
						randomRule(),
						randomRule(),
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
}
