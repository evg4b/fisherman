package rules

import (
	"fisherman/internal"
	"fisherman/utils"
	"io"
	"io/fs"

	"github.com/spf13/afero"
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

	return afero.WriteFile(ctx.Files(), args[0], []byte(rule.Message), fs.ModePerm)
}

func (rule *PrepareMessage) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Message, variables)
}
