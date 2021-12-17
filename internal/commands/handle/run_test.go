package handle_test

import (
	"context"
	"errors"
	. "fisherman/internal/commands/handle"
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
	"fisherman/internal/rules"
	"fisherman/pkg/log"
	"fisherman/pkg/vcs"
	"fisherman/testing/mocks"
	"io"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var globalVars = map[string]interface{}{
	constants.BranchNameVariable:       "/refs/head/develop",
	constants.TagVariable:              "1.0.0",
	constants.UserEmailVariable:        "evg4b@mail.com",
	constants.UserNameVariable:         "evg4b",
	constants.FishermanVersionVariable: constants.Version,
	constants.CwdVariable:              "~/project",
	constants.OsVariable:               runtime.GOOS,
}

func TestCommand_Run(t *testing.T) {
	log.SetOutput(io.Discard)

	repoStub := mocks.NewRepositoryMock(t).
		GetCurrentBranchMock.Return("/refs/head/develop", nil).
		GetLastTagMock.Return("1.0.0", nil).
		GetUserMock.Return(vcs.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil)

	validRule := mocks.NewRuleMock(t).
		GetTypeMock.Return(rules.ExecType).
		ConfigureMock.Return().
		GetPositionMock.Return(rules.Scripts).
		CompileMock.Return().
		GetContitionMock.Return("").
		GetPrefixMock.Return("prefix-").
		CheckMock.Return(nil)

	t.Run("runs correctly", func(t *testing.T) {
		command := NewCommand(
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: []configuration.Rule{validRule},
				},
			}),
			WithRepository(repoStub),
			WithWorkersCount(5),
		)

		err := command.Run(context.TODO(), []string{"--hook", "pre-commit"})

		assert.NoError(t, err)
	})

	t.Run("unknown hook", func(t *testing.T) {
		command := NewCommand(
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: []configuration.Rule{
						mocks.NewRuleMock(t),
					},
				},
			}),
			WithRepository(repoStub),
		)

		err := command.Run(context.TODO(), []string{"--hook", "test"})

		assert.EqualError(t, err, "'test' is not valid hook name")
	})

	t.Run("call handler and return error", func(t *testing.T) {
		command := NewCommand(
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: []configuration.Rule{
						validRule.CheckMock.Return(errors.New("test error")),
					},
				},
			}),
			WithRepository(repoStub),
		)

		err := command.Run(context.TODO(), []string{"--hook", "pre-commit"})

		assert.EqualError(t, err, "1 error occurred:\n\t* [exec] test error\n\n")
	})

	t.Run("call handler with global variables", func(t *testing.T) {
		command := NewCommand(
			WithHooksConfig(&configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: []configuration.Rule{validRule},
				},
			}),
			WithGlobalVars(globalVars),
			WithRepository(mocks.NewRepositoryMock(t).
				GetCurrentBranchMock.Return("/refs/head/develop", nil).
				GetLastTagMock.Return("1.0.0", errors.New("test error")).
				GetUserMock.Return(vcs.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil)),
		)

		err := command.Run(context.TODO(), []string{"--hook", "pre-commit"})

		assert.EqualError(t, err, "1 error occurred:\n\t* [exec] test error\n\n")
	})
}
