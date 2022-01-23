// nolint: dupl
package rules_test

import (
	"context"
	. "fisherman/internal/rules"
	"fisherman/testing/testutils"
	"io/ioutil"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
)

func TestCommitMessage_Compile(t *testing.T) {
	rule := CommitMessage{
		BaseRule: BaseRule{Type: CommitMessageType},
		Prefix:   "Prefix{{var1}}",
		Suffix:   "Suffix{{var1}}",
		Regexp:   "Regexp{{var1}}",
		NotEmpty: true,
	}

	rule.Compile(map[string]interface{}{"var1": "VALUE"})

	assert.Equal(t, CommitMessage{
		BaseRule: BaseRule{Type: CommitMessageType},
		Prefix:   "PrefixVALUE",
		Suffix:   "SuffixVALUE",
		Regexp:   "RegexpVALUE",
		NotEmpty: true,
	}, rule)
}

func TestCommitMessage_Check(t *testing.T) {
	t.Run("not-empty field", func(t *testing.T) {
		tests := []struct {
			name        string
			message     string
			notEmpty    bool
			expectedErr string
		}{
			{
				name:        "Active with empty string",
				notEmpty:    true,
				expectedErr: "[commit-message] commit message should not be empty",
			},
			{
				name:     "Inactive with empty string",
				notEmpty: false,
			},
			{
				name:     "Active with not empty string",
				message:  "message",
				notEmpty: true,
			},
			{
				name:     "Active with not empty string",
				message:  "message",
				notEmpty: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rule := makeRule(
					&CommitMessage{
						BaseRule: BaseRule{Type: CommitMessageType},
						NotEmpty: tt.notEmpty,
					},
					WithArgs([]string{"massage"}),
					WithFileSystem(testutils.FsFromMap(t, map[string]string{
						"massage": tt.message,
					})),
				)

				err := rule.Check(context.TODO(), ioutil.Discard)

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("prefix field", func(t *testing.T) {
		tests := []struct {
			name        string
			message     string
			prefix      string
			expectedErr string
		}{
			{
				name:        "active with empty string",
				expectedErr: "[commit-message] commit message should have prefix 'prefix-'",
				prefix:      "prefix-",
			},
			{
				name:   "inactive with empty string",
				prefix: "",
			},
			{
				name:    "active with string and prefix",
				message: "prefix-message",
				prefix:  "prefix-",
			},
			{
				name:    "inactive with string and prefix",
				message: "prefix-message",
			},
			{
				name:        "active with string and other prefix",
				expectedErr: "[commit-message] commit message should have prefix 'prefix-'",
				message:     "other-prefix-message",
				prefix:      "prefix-",
			},
			{
				name:    "inactive with string and other prefix",
				message: "other-prefix-message",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rule := makeRule(
					&CommitMessage{
						BaseRule: BaseRule{Type: CommitMessageType},
						Prefix:   tt.prefix,
					},
					WithArgs([]string{"massage"}),
					WithFileSystem(testutils.FsFromMap(t, map[string]string{
						"massage": tt.message,
					})),
				)

				err := rule.Check(context.TODO(), ioutil.Discard)

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("suffix field", func(t *testing.T) {
		tests := []struct {
			name        string
			message     string
			suffix      string
			expectedErr string
		}{
			{
				name:        "active with empty string",
				expectedErr: "[commit-message] commit message should have suffix '-suffix'",
				message:     "",
				suffix:      "-suffix",
			},
			{
				name:   "inactive with empty string",
				suffix: "",
			},
			{
				name:    "active with string and suffix",
				message: "message-suffix",
				suffix:  "-suffix",
			},
			{
				name:    "inactive with string and suffix",
				message: "message-suffix",
				suffix:  "",
			},
			{
				name:        "active with string and other suffix",
				expectedErr: "[commit-message] commit message should have suffix '-suffix'",
				message:     "message-suffix-other",
				suffix:      "-suffix",
			},
			{
				name:    "inactive with string and other suffix",
				message: "message-suffix-other",
				suffix:  "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rule := makeRule(
					&CommitMessage{
						BaseRule: BaseRule{Type: CommitMessageType},
						Suffix:   tt.suffix,
					},
					WithArgs([]string{"massage"}),
					WithFileSystem(testutils.FsFromMap(t, map[string]string{
						"massage": tt.message,
					})),
				)

				err := rule.Check(context.TODO(), ioutil.Discard)

				testutils.AssertError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("regexp field", func(t *testing.T) {
		t.Run("correct mattching", func(t *testing.T) {
			tests := []struct {
				name        string
				message     string
				expression  string
				expectedErr string
			}{
				{
					name:       "inactive with empty string",
					expression: "",
				},
				{
					name:        "active with empty string",
					expectedErr: "[commit-message] commit message should be matched regular expression '^[a-z]{5}$'",
					expression:  "^[a-z]{5}$",
				},
				{
					name:        "active with correct matching",
					expectedErr: "[commit-message] commit message should be matched regular expression '^[a-z]{5}$'",
					message:     "message",
					expression:  "^[a-z]{5}$",
				},
				{
					name:       "active with correct matching",
					message:    "message",
					expression: "^[a-z]{7}$",
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					rule := makeRule(
						&CommitMessage{
							BaseRule: BaseRule{Type: CommitMessageType},
							Regexp:   tt.expression,
						},
						WithArgs([]string{"massage"}),
						WithFileSystem(testutils.FsFromMap(t, map[string]string{
							"massage": tt.message,
						})),
					)

					err := rule.Check(context.TODO(), ioutil.Discard)

					testutils.AssertError(t, tt.expectedErr, err)
				})
			}
		})

		t.Run("regexp parsing error", func(t *testing.T) {
			rule := CommitMessage{
				BaseRule: BaseRule{Type: CommitMessageType},
				Regexp:   "[a-z]($",
			}

			err := rule.Check(context.TODO(), ioutil.Discard)

			assert.Error(t, err)
		})
	})

	t.Run("message reading error", func(t *testing.T) {
		rule := makeRule(
			&CommitMessage{
				BaseRule: BaseRule{Type: CommitMessageType},
				Regexp:   "[a-z]($",
			},
			WithArgs([]string{"unknow/file"}),
			WithFileSystem(memfs.New()),
		)

		err := rule.Check(context.TODO(), ioutil.Discard)

		assert.EqualError(t, err, "message cannot be read: file does not exist")
	})
}
