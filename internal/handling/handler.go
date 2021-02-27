package handling

import (
	"fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/expression"
)

type Handler interface {
	Handle(ctx internal.ExecutionContext, args []string) error
}

type Action = func(internal.ExecutionContext) (bool, error)

type HookHandler struct {
	Engine          expression.Engine
	Rules           []configuration.Rule
	Scripts         []configuration.Rule
	PostScriptRules []configuration.Rule
	WorkersCount    int
}

func (handler *HookHandler) Handle(ctx internal.ExecutionContext, args []string) error {
	err := handler.runRules(ctx, handler.Rules)
	if err != nil {
		return err
	}

	err = handler.runRules(ctx, handler.Scripts)
	if err != nil {
		return err
	}

	return handler.runRules(ctx, handler.PostScriptRules)
}
