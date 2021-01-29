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

func (config PrepareMessage) Check(_ io.Writer, ctx internal.ExecutionContext) error {
	if utils.IsEmpty(config.Message) {
		return nil
	}

	args := ctx.Args()
	files := ctx.Files()

	return files.Write(args[0], config.Message)
}
