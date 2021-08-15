package rules

import (
	"fisherman/internal"
	"fisherman/internal/utils"
	"io"
	"io/fs"

	"github.com/go-git/go-billy/v5/util"
)

const PrepareMessageType = "prepare-message"

type PrepareMessage struct {
	BaseRule `yaml:",inline"`
	Message  string `yaml:"message"`
}

func (rule PrepareMessage) Check(ctx internal.ExecutionContext, _ io.Writer) error {
	if utils.IsEmpty(rule.Message) {
		return nil
	}

	args := ctx.Args()

	return util.WriteFile(ctx.Files(), args[0], []byte(rule.Message), fs.ModePerm)
}

func (rule *PrepareMessage) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Message, variables)
}
