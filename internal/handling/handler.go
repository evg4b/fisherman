package handling

import (
	"fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/expression"
	"fisherman/internal/validation"
)

type Handler interface {
	Handle(ctx internal.AsyncContext, args []string) error
}

type Action = func(internal.SyncContext) (bool, error)

type HookHandler struct {
	Engine          expression.Engine
	BeforeActions   []Action
	Rules           []configuration.Rule
	Scripts         configuration.ScriptsConfig
	AsyncValidators []validation.AsyncValidator
	PostScriptRules []configuration.Rule
	AfterActions    []Action
	WorkersCount    int
}

func (handler *HookHandler) Handle(ctx internal.AsyncContext, args []string) error {
	next, err := RunActions(ctx, handler.BeforeActions)
	if err != nil || !next {
		return err
	}

	err = handler.runRules(ctx, handler.Rules)
	if err != nil {
		return err
	}

	err = validation.RunAsync(ctx, handler.AsyncValidators)
	if err != nil {
		return err
	}

	err = handler.runRules(ctx, handler.PostScriptRules)
	if err != nil {
		return err
	}

	_, err = RunActions(ctx, handler.AfterActions)

	return err
}

func RunActions(ctx internal.SyncContext, actions []Action) (bool, error) {
	for _, action := range actions {
		next, err := action(ctx)
		if err != nil || !next {
			return false, err
		}
	}

	return true, nil
}
