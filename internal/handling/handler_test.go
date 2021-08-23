package handling_test

import (
	"fisherman/internal/configuration"
	"fisherman/internal/handling"
	"fisherman/internal/validation"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"io"
	"testing"
)

// TODO: Add test for ctx Cancel

// nolint: dupl
func TestHookHandler_Handle_Rules(t *testing.T) {
	ctx := mocks.NewExecutionContextMock(t).
		OutputMock.Return(testutils.NopCloser(io.Discard)).
		GlobalVariablesMock.Return(map[string]interface{}{}, nil)

	tests := []struct {
		name        string
		handler     handling.HookHandler
		expectedErr string
	}{
		{
			name: "positive case",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				Rules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
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
					EvalMock.Return(false, validation.Errorf("test-rule", "test")),
				Rules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1"),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2"),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3"),
				},
				WorkersCount: 2,
			},
			expectedErr: "[test-rule] test",
		},
		{
			name: "rule returns error",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				Rules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3").
						CheckMock.Return(validation.Errorf("test-rule", "test")),
				},
				WorkersCount: 2,
			},
			expectedErr: "1 error occurred:\n\t* [rule3] [test-rule] test\n\n",
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
			expectedErr: "incorrect workers count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.handler.Handle(ctx, []string{})

			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}

// nolint: dupl
func TestHookHandler_Handle_PostScriptRules(t *testing.T) {
	ctx := mocks.NewExecutionContextMock(t).
		OutputMock.Return(testutils.NopCloser(io.Discard)).
		GlobalVariablesMock.Return(map[string]interface{}{}, nil)

	tests := []struct {
		name        string
		handler     handling.HookHandler
		expectedErr string
	}{
		{
			name: "positive case",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
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
					EvalMock.Return(false, validation.Errorf("test-rule", "test")),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1"),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2"),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3"),
				},
				WorkersCount: 2,
			},
			expectedErr: "[test-rule] test",
		},
		{
			name: "rule returns error",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t).
					EvalMock.Return(true, nil),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==0").
						GetTypeMock.Return("rule1").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==2").
						GetTypeMock.Return("rule2").
						CheckMock.Return(nil),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Expect().Return("1==3").
						GetTypeMock.Return("rule3").
						CheckMock.Return(validation.Errorf("test-rule", "test")),
				},
				WorkersCount: 2,
			},
			expectedErr: "1 error occurred:\n\t* [rule3] [test-rule] test\n\n",
		},
		{
			name: "incorrect workers count",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t),
					mocks.NewRuleMock(t),
					mocks.NewRuleMock(t),
				},
				WorkersCount: 0,
			},
			expectedErr: "incorrect workers count",
		},
		{
			name: "empty condition",
			handler: handling.HookHandler{
				Engine: mocks.NewEngineMock(t),
				PostScriptRules: []configuration.Rule{
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Return("").
						CheckMock.Return(nil).
						GetTypeMock.Return("test1"),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Return("").
						CheckMock.Return(nil).
						GetTypeMock.Return("test2"),
					mocks.NewRuleMock(t).
						GetPrefixMock.Return("TEST").
						GetContitionMock.Return("").
						CheckMock.Return(nil).
						GetTypeMock.Return("test3"),
				},
				WorkersCount: 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.handler.Handle(ctx, []string{})

			testutils.CheckError(t, tt.expectedErr, err)
			print(tt.name)
		})
	}
}
