package handlers

import (
	"errors"
	"fisherman/config/hooks"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestValidateMessageNotEmpty(t *testing.T) {
	err := errors.New("Commit comment should not be empty")
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
			assertMultiError(t, validateMessage(tt.message, &hooks.CommitMsgHookConfig{NotEmpty: tt.notEmpty}), tt.err)
		})
	}
}

func TestValidateMessageCommitPrefix(t *testing.T) {
	err := errors.New("Commit should have prefix '[prefix]'")
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
			assertMultiError(t, validateMessage(tt.message, &hooks.CommitMsgHookConfig{CommitPrefix: "[prefix]"}), tt.err)
		})
	}
}

func TestValidateMessageCommitSuffix(t *testing.T) {
	err := errors.New("Commit should have suffix '[suffix]'")
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
			assertMultiError(t, validateMessage(tt.message, &hooks.CommitMsgHookConfig{CommitSuffix: "[suffix]"}), tt.err)
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
