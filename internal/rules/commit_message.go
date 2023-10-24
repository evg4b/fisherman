package rules

import (
	"context"
	"fmt"
	"github.com/evg4b/fisherman/internal/utils"
	"io"
	"regexp"
	"strings"

	"github.com/go-git/go-billy/v5/util"
)

const CommitMessageType = "commit-message"

const filePathArgumentIndex = 0

type CommitMessage struct {
	BaseRule `yaml:",inline"`
	Prefix   string `yaml:"prefix"`
	Suffix   string `yaml:"suffix"`
	Regexp   string `yaml:"regexp"`
	NotEmpty bool   `yaml:"not-empty"`
}

// nolint: cyclop
func (rule CommitMessage) Check(_ context.Context, _ io.Writer) error {
	message, err := rule.Message()
	if err != nil {
		return err
	}

	if rule.NotEmpty && utils.IsEmpty(message) {
		return rule.errorf("commit message should not be empty")
	}

	if !utils.IsEmpty(rule.Prefix) && !strings.HasPrefix(message, rule.Prefix) {
		return rule.errorf("commit message should have prefix '%s'", rule.Prefix)
	}

	if !utils.IsEmpty(rule.Suffix) && !strings.HasSuffix(message, rule.Suffix) {
		return rule.errorf("commit message should have suffix '%s'", rule.Suffix)
	}

	if !utils.IsEmpty(rule.Regexp) {
		matched, err := regexp.MatchString(rule.Regexp, message)
		if err != nil {
			return err
		}

		if !matched {
			return rule.errorf("commit message should be matched regular expression '%s'", rule.Regexp)
		}
	}

	return nil
}

func (rule *CommitMessage) Compile(variables map[string]any) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Prefix, variables)
	utils.FillTemplate(&rule.Suffix, variables)
	utils.FillTemplate(&rule.Regexp, variables)
}

func (rule *CommitMessage) Message() (string, error) {
	messageFilePath, err := rule.arg(filePathArgumentIndex)
	if err != nil {
		return "", err
	}

	message, err := util.ReadFile(rule.BaseRule.fs, messageFilePath)
	if err != nil {
		return "", fmt.Errorf("message cannot be read: %w", err)
	}

	return string(message), nil
}
