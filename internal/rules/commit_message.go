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
	BaseRule `mapstructure:",squash"`
	Prefix   string `mapstructure:"prefix"`
	Suffix   string `mapstructure:"suffix"`
	Regexp   string `mapstructure:"regexp"`
	NotEmpty bool   `mapstructure:"not-empty"`
}

func (rule CommitMessage) Check(ctx internal.ExecutionContext, _ io.Writer) error {
	message := ctx.Message()

	if rule.NotEmpty && utils.IsEmpty(message) {
		return fmt.Errorf("commit message should not be empty")
	}

	if !utils.IsEmpty(rule.Prefix) && !strings.HasPrefix(ctx.Message(), rule.Prefix) {
		return fmt.Errorf("commit message should have prefix '%s'", rule.Prefix)
	}

	if !utils.IsEmpty(rule.Suffix) && !strings.HasSuffix(ctx.Message(), rule.Suffix) {
		return fmt.Errorf("commit message should have suffix '%s'", rule.Suffix)
	}

	if !utils.IsEmpty(rule.Regexp) {
		matched, err := regexp.MatchString(rule.Regexp, ctx.Message())
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
