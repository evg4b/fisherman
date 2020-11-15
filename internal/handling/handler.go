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
	contextFactory  internal.CtxFactory
	beforeActions   []Action
	syncValidators  []validation.SyncValidator
	asyncValidators []validation.AsyncValidator
	afterActions    []Action
}

func NewHookHandler(
	contextFactory internal.CtxFactory,
	beforeActions []Action,
	syncValidators []validation.SyncValidator,
	asyncValidators []validation.AsyncValidator,
	afterActions []Action,
) *HookHandler {
	return &HookHandler{
		beforeActions:   beforeActions,
		contextFactory:  contextFactory,
		asyncValidators: asyncValidators,
		syncValidators:  syncValidators,
		afterActions:    afterActions,
	}
}

func (h *HookHandler) Handle(args []string) error {
	ctx := h.contextFactory(args, log.InfoOutput)
	next, err := RunActions(ctx, h.beforeActions)
	if err != nil || !next {
		return err
	}

	err = validation.RunSync(ctx, h.syncValidators)
	if err != nil {
		return err
	}

	err = validation.RunAsync(ctx, h.asyncValidators)
	if err != nil {
		return err
	}

	_, err = RunActions(ctx, h.afterActions)

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
