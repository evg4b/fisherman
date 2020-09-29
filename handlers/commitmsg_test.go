package handlers

import (
	"errors"
	"fisherman/config/hooks"
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
