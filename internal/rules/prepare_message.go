package rules

import (
	"fisherman/internal"
	"fisherman/utils"
	"io"
)

const PrepareMessageType = "prepare-message"

type PrepareMessage struct {
	BaseRule `mapstructure:",squash"`
	Message  string `mapstructure:"message"`
}

func (rule PrepareMessage) Check(ctx internal.ExecutionContext, _ io.Writer) error {
	if utils.IsEmpty(rule.Message) {
		return nil
	}

	args := ctx.Args()

	return ctx.Files().Write(args[0], rule.Message)
}

func (rule *PrepareMessage) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Message, variables)
}
