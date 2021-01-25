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

func (config CommitMessage) Check(_ io.Writer, ctx internal.ExecutionContext) error {
	message := ctx.Message()

	if config.NotEmpty && utils.IsEmpty(message) {
		return fmt.Errorf("commit message should not be empty")
	}

	if !utils.IsEmpty(config.Prefix) && !strings.HasPrefix(ctx.Message(), config.Prefix) {
		return fmt.Errorf("commit message should have prefix '%s'", config.Prefix)
	}

	if !utils.IsEmpty(config.Suffix) && !strings.HasSuffix(ctx.Message(), config.Suffix) {
		return fmt.Errorf("commit message should have suffix '%s'", config.Suffix)
	}

	if !utils.IsEmpty(config.Regexp) {
		matched, err := regexp.MatchString(config.Regexp, ctx.Message())
		if err != nil {
			return err
		}

		if !matched {
			return fmt.Errorf("commit message should be matched regular expression '%s'", config.Regexp)
		}
	}

	return nil
}
