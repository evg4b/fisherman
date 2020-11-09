package handlers

import (
	"context"
	"errors"
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/config/hooks"
	"fisherman/infrastructure"
	inf_mock "fisherman/mocks/infrastructure"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestValidateMessageNotEmpty(t *testing.T) {
	err := errors.New("commit message should not be empty")
	testData := []struct {
		message  string
		notEmpty bool
		config   hooks.CommitMsgHookConfig
		err      error
	}{
		{message: "not empty string", notEmpty: true, err: nil},
		{message: "  not empty string", notEmpty: true, err: nil},
		{message: "", notEmpty: true, err: err},
		{message: "", notEmpty: false, err: nil},
		{message: "   ", notEmpty: true, err: err},
		{message: "   ", notEmpty: false, err: nil},
	}

	for _, tt := range testData {
		t.Run(tt.message, func(t *testing.T) {
			actualError := validateMessage(tt.message, &hooks.CommitMsgHookConfig{NotEmpty: tt.notEmpty})
			assertMultiError(t, actualError, tt.err)
		})
	}
}

func TestValidateMessageCommitPrefix(t *testing.T) {
	err := errors.New("commit message should have prefix '[prefix]'")
	config := hooks.CommitMsgHookConfig{MessagePrefix: "[prefix]"}

	testData := []struct {
		message string
		err     error
	}{
		{message: "[prefix] message", err: nil},
		{message: "message", err: err},
		{message: " [prefix] message", err: err},
		{message: "message[prefix]", err: err},
	}

	for _, tt := range testData {
		t.Run(tt.message, func(t *testing.T) {
			actualError := validateMessage(tt.message, &config)
			assertMultiError(t, actualError, tt.err)
		})
	}
}

func TestValidateMessageCommitSuffix(t *testing.T) {
	err := errors.New("commit message should have suffix '[suffix]'")
	config := hooks.CommitMsgHookConfig{MessageSuffix: "[suffix]"}

	testData := []struct {
		message string
		err     error
	}{
		{message: "[suffix] message", err: err},
		{message: "message", err: err},
		{message: "message [suffix] ", err: err},
		{message: "message [suffix]", err: nil},
	}

	for _, tt := range testData {
		t.Run(tt.message, func(t *testing.T) {
			actualError := validateMessage(tt.message, &config)
			assertMultiError(t, actualError, tt.err)
		})
	}
}

func TestValidateMessageCommitRegexp(t *testing.T) {
	testData := []struct {
		message string
		regexp  string
		err     error
	}{
		{message: "message", regexp: "", err: nil},
		{message: "message", regexp: "^[a-z]*$", err: nil},
		{
			message: "Message",
			regexp:  "^[a-z]*$",
			err:     errors.New("commit message should be matched regular expression '^[a-z]*$'"),
		},
	}

	for _, tt := range testData {
		t.Run(tt.message, func(t *testing.T) {
			config := hooks.CommitMsgHookConfig{MessageRegexp: tt.regexp}
			actualError := validateMessage(tt.message, &config)
			assertMultiError(t, actualError, tt.err)
		})
	}
}

func assertMultiError(t *testing.T, multipleErrors *multierror.Error, expectedError error) {
	if expectedError != nil {
		assert.NotNil(t, multipleErrors)
		assert.Contains(t, multipleErrors.Errors, expectedError)
	} else {
		assert.Nil(t, multipleErrors)
	}
}

func TestCommitMsgHandler(t *testing.T) {
	fakeRepo := inf_mock.Repository{}
	fakeRepo.On("GetCurrentBranch").Return("develop", nil)
	fakeRepo.On("GetLastTag").Return("0.0.0", nil)
	fakeRepo.On("GetUser").Return(infrastructure.User{}, nil)

	fakeFS := inf_mock.FileSystem{}
	fakeFS.On("Read", ".git/MESSAGE").Return("[fisherman] test commit", nil)

	tests := []struct {
		name string
		args []string
		err  error
	}{
		{name: "base test", args: []string{".git/MESSAGE"}, err: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
				Config:     &config.DefaultConfig,
				Repository: &fakeRepo,
				FileSystem: &fakeFS,
				App:        &clicontext.AppInfo{},
			})
			err := CommitMsgHandler(ctx, tt.args)
			assert.Equal(t, tt.err, err)
		})
	}
}
