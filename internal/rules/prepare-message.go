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

func (config PrepareMessage) Check(ctx internal.ExecutionContext, _ io.Writer) error {
	if utils.IsEmpty(config.Message) {
		return nil
	}

	args := ctx.Args()

	return ctx.Files().Write(args[0], config.Message)
}
