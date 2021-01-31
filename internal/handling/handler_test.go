package handling_test

import (
	"errors"
	"fisherman/configuration"
	"fisherman/internal/handling"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"testing"
)

// nolint: dupl
func TestHookHandler_Handle_Rules(t *testing.T) {
	ctx := mocks.NewExecutionContextMock(t)

	tests := []struct {
		name    string
		handler handling.HookHandler
		err     string
	}{
		{
			name: "positive case",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				Rules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3").
						CheckMock.Return(nil),
				},
				WorkersCount: 2,
			},
		},
		{
			name: "negative case",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(false, errors.New("test")),
				Rules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1"),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2"),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3"),
				},
				WorkersCount: 2,
			},
			err: "test",
		},
		{
			name: "rule returns error",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				Rules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3").
						CheckMock.Return(errors.New("test")),
				},
				WorkersCount: 2,
			},
			err: "1 error occurred:\n\t* [rule3] test\n\n",
		},
		{
			name: "rule returns error",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				Rules: []configuration.Rule{
					mocks.NewRuleMock(t),
					mocks.NewRuleMock(t),
					mocks.NewRuleMock(t),
				},
				WorkersCount: 0,
			},
			err: "incorrect workers count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.handler.Handle(ctx, []string{})

			testutils.CheckError(t, tt.err, err)
		})
	}
}

// nolint: dupl
func TestHookHandler_Handle_PostScriptRules(t *testing.T) {
	ctx := mocks.NewExecutionContextMock(t)

	tests := []struct {
		name    string
		handler handling.HookHandler
		err     string
	}{
		{
			name: "positive case",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3").
						CheckMock.Return(nil),
				},
				WorkersCount: 2,
			},
		},
		{
			name: "negative case",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(false, errors.New("test")),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1"),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2"),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3"),
				},
				WorkersCount: 2,
			},
			err: "test",
		},
		{
			name: "rule returns error",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3").
						CheckMock.Return(errors.New("test")),
				},
				WorkersCount: 2,
			},
			err: "1 error occurred:\n\t* [rule3] test\n\n",
		},
		{
			name: "rule returns error",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t),
					mocks.NewRuleMock(t),
					mocks.NewRuleMock(t),
				},
				WorkersCount: 0,
			},
			err: "incorrect workers count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.handler.Handle(ctx, []string{})

			testutils.CheckError(t, tt.err, err)
		})
	}
}
