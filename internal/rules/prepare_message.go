package rules

import (
	"context"
	"github.com/evg4b/fisherman/internal/utils"
	"io"
	"io/fs"

	"github.com/go-git/go-billy/v5/util"
)

const PrepareMessageType = "prepare-message"

type PrepareMessage struct {
	BaseRule `yaml:",inline"`
	Message  string `yaml:"message"`
}

func (rule *PrepareMessage) Check(_ context.Context, _ io.Writer) error {
	if utils.IsEmpty(rule.Message) {
		return nil
	}

	arg, err := rule.arg(0)
	if err != nil {
		return err
	}

	return util.WriteFile(rule.BaseRule.fs, arg, []byte(rule.Message), fs.ModePerm)
}

func (rule *PrepareMessage) Compile(variables map[string]any) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Message, variables)
}
