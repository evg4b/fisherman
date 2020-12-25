package handling

import (
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"fisherman/internal/validation"
)

type Handler interface {
	Handle(args []string) error
}

type Action = func(internal.SyncContext) (bool, error)

type HookHandler struct {
	ContextFactory  internal.CtxFactory
	BeforeActions   []Action
	SyncValidators  []validation.SyncValidator
	AsyncValidators []validation.AsyncValidator
	AfterActions    []Action
}

func (h *HookHandler) Handle(args []string) error {
	ctx := h.ContextFactory(args, log.InfoOutput)
	next, err := RunActions(ctx, h.BeforeActions)
	if err != nil || !next {
		return err
	}

	err = validation.RunSync(ctx, h.SyncValidators)
	if err != nil {
		return err
	}

	err = validation.RunAsync(ctx, h.AsyncValidators)
	if err != nil {
		return err
	}

	_, err = RunActions(ctx, h.AfterActions)

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
