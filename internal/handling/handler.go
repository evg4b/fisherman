package handling

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
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
	GlobalVariables Variables
	WorkersCount    int
}

func (h *HookHandler) Handle(ctx internal.ExecutionContext, args []string) error {
	err := h.runRules(ctx, h.Rules)
	if err != nil {
		return err
	}

	err = h.runRules(ctx, h.Scripts)
	if err != nil {
		return err
	}

	return h.runRules(ctx, h.PostScriptRules)
}
