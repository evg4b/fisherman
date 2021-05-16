package rules

import (
	"fisherman/internal"
	"fisherman/utils"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const CommitMessageType = "commit-message"

type CommitMessage struct {
	BaseRule `yaml:",inline"`
	Prefix   string `yaml:"prefix"`
	Suffix   string `yaml:"suffix"`
	Regexp   string `yaml:"regexp"`
	NotEmpty bool   `yaml:"not-empty"`
}

func (rule CommitMessage) Check(ctx internal.ExecutionContext, _ io.Writer) error {
	message, err := ctx.Message()
	if err != nil {
		return err
	}

	if rule.NotEmpty && utils.IsEmpty(message) {
		return fmt.Errorf("commit message should not be empty")
	}

	if !utils.IsEmpty(rule.Prefix) && !strings.HasPrefix(message, rule.Prefix) {
		return fmt.Errorf("commit message should have prefix '%s'", rule.Prefix)
	}

	if !utils.IsEmpty(rule.Suffix) && !strings.HasSuffix(message, rule.Suffix) {
		return fmt.Errorf("commit message should have suffix '%s'", rule.Suffix)
	}

	if !utils.IsEmpty(rule.Regexp) {
		matched, err := regexp.MatchString(rule.Regexp, message)
		if err != nil {
			return err
		}

		if !matched {
			return fmt.Errorf("commit message should be matched regular expression '%s'", rule.Regexp)
		}
	}

	return nil
}

func (rule *CommitMessage) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Prefix, variables)
	utils.FillTemplate(&rule.Suffix, variables)
	utils.FillTemplate(&rule.Regexp, variables)
}
