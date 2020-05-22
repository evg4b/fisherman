package handling_test

import (
	"context"
	"errors"
	"github.com/evg4b/fisherman/internal"
	"github.com/evg4b/fisherman/internal/configuration"
	"github.com/evg4b/fisherman/internal/constants"
	"github.com/evg4b/fisherman/internal/expression"
	"github.com/evg4b/fisherman/internal/rules"
	"github.com/evg4b/fisherman/internal/validation"
	"github.com/evg4b/fisherman/pkg/guards"
	"github.com/evg4b/fisherman/pkg/vcs"
	"github.com/evg4b/fisherman/testing/mocks"
	"github.com/evg4b/fisherman/testing/testutils"
	"io"
	"testing"

	. "github.com/evg4b/fisherman/internal/handling"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHookHandler_Handle(t *testing.T) {
	engine := expression.NewGoExpressionEngine()

	validRule := mocks.NewRuleMock(t).
		GetTypeMock.Return(rules.ExecType).
		ConfigureMock.Return().
		GetPositionMock.Return(rules.Scripts).
		CompileMock.Return().
		GetContitionMock.Return("").
		GetPrefixMock.Return("prefix-").
		CheckMock.Return(nil)

	repo := mocks.NewRepositoryMock(t).
		GetUserMock.Return(vcs.User{UserName: "", Email: ""}, nil).
		GetCurrentBranchMock.Return("main", nil).
		GetLastTagMock.Return("0.0.1", nil)

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
					ConfigureMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST2", "1 == 1", rules.ExecType).
					ConfigureMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST3", "1 == 3", rules.SuppressCommitFilesType).
					ConfigureMock.Return().
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
					ConfigureMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST", "1 == 1", rules.ExecType).
					ConfigureMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST", "1 == 2", rules.SuppressCommitFilesType).
					ConfigureMock.Return().
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
					ConfigureMock.Return().
					CheckMock.Return(nil),
				ruleMock(t, "TEST", "1 == 1", rules.ExecType).
					ConfigureMock.Return().
					CheckMock.Return(validation.Errorf("test-rule", "test")),
				ruleMock(t, "TEST", "1 == 3", rules.SuppressCommitFilesType).
					ConfigureMock.Return().
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
					WithRepository(repo),
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
					WithRepository(repo),
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
					WithRepository(repo),
				)

				guards.NoError(err)

				err = handler.Handle(context.TODO())

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("run rules in correct orders", func(t *testing.T) {
		var order []string

		ruleDef := func(ruleType string, position byte) configuration.Rule {
			checkHandler := func() func(context.Context, io.Writer) error {
				return func(ec context.Context, w io.Writer) error {
					order = append(order, ruleType)

					return nil
				}
			}

			return mocks.NewRuleMock(t).
				GetTypeMock.Return(rules.ExecType).
				ConfigureMock.Return().
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
			WithRepository(repo),
		)

		guards.NoError(err)

		err = handler.Handle(context.TODO())

		require.NoError(t, err)
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
			WithRepository(repo),
		)

		guards.NoError(err)

		ctx, cancel := context.WithCancel(context.Background())

		cancel()

		err = handler.Handle(ctx)

		require.EqualError(t, err, "1 error occurred:\n\t* context canceled\n\n")
	})

	t.Run("hook config is not presented", func(t *testing.T) {
		_, err := NewHookHandler(
			constants.PreCommitHook,
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: nil,
			}),
		)

		require.Error(t, err, ErrNotPresented.Error())
	})

	t.Run("invalid hook name", func(t *testing.T) {
		_, err := NewHookHandler(
			"unknown-hook",
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: nil,
			}),
		)

		require.Error(t, err, "'unknown-hook' is not valid hook name")
	})

	t.Run("fails when get predefined variables", func(t *testing.T) {
		predefinedVarsTestCases2 := []struct {
			name        string
			repo        internal.Repository
			expectedErr string
		}{
			{
				name: "get last tag error",
				repo: mocks.NewRepositoryMock(t).
					GetCurrentBranchMock.Return("/refs/head/develop", nil).
					GetLastTagMock.Return("", errors.New("tag getting error")).
					GetUserMock.Return(vcs.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil),
				expectedErr: "tag getting error",
			},
			{
				name: "get current branch error",
				repo: mocks.NewRepositoryMock(t).
					GetCurrentBranchMock.Return("", errors.New("branch getting error")).
					GetLastTagMock.Return("1.0.0", nil).
					GetUserMock.Return(vcs.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil),
				expectedErr: "branch getting error",
			},
			{
				name: "get git user error",
				repo: mocks.NewRepositoryMock(t).
					GetCurrentBranchMock.Return("/refs/head/develop", nil).
					GetLastTagMock.Return("1.0.0", nil).
					GetUserMock.Return(vcs.User{}, errors.New("user getting error")),
				expectedErr: "user getting error",
			},
		}

		for _, tt := range predefinedVarsTestCases2 {
			t.Run(tt.name, func(t *testing.T) {
				handler, err := NewHookHandler(
					"pre-commit",
					WithExpressionEngine(mocks.NewEngineMock(t)),
					WithHooksConfig(&configuration.HooksConfig{
						PreCommitHook: &configuration.HookConfig{
							Rules: []configuration.Rule{validRule},
						},
					}),
					WithWorkersCount(4),
					WithFileSystem(mocks.NewFilesystemMock(t)),
					WithRepository(tt.repo),
				)

				assert.Nil(t, handler)
				require.EqualError(t, err, tt.expectedErr)
			})
		}
	})
}

func ruleMock(t *testing.T, prefix, condition, ruleType string) *mocks.RuleMock {
	t.Helper()

	return mocks.NewRuleMock(t).
		GetPrefixMock.Return(prefix).
		GetContitionMock.Return(condition).
		GetTypeMock.Return(ruleType).
		CompileMock.Return().
		ConfigureMock.Return()
}

func applyPosition(rules []*mocks.RuleMock, position byte) []configuration.Rule {
	updatedRules := make([]configuration.Rule, 0, len(rules))
	for _, rule := range rules {
		rule.GetPositionMock.Return(position)
		updatedRules = append(updatedRules, rule)
	}

	return updatedRules
}
