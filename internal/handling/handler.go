package handling

import (
	c "fisherman/config"
	"fisherman/infrastructure/log"
	i "fisherman/internal"
	v "fisherman/internal/validation"
)

type Handler interface {
	IsConfigured(config *c.HooksConfig) bool
	Handle(args []string) error
}

type Action = func(i.SyncContext) (bool, error)

type HookHandler struct {
	contextFactory  i.CtxFactory
	beforeActions   []Action
	syncValidators  []v.SyncValidator
	asyncValidators []v.AsyncValidator
	afterActions    []Action
}

func NewHookHandler(
	contextFactory i.CtxFactory,
	beforeActions []Action,
	syncValidators []v.SyncValidator,
	asyncValidators []v.AsyncValidator,
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

	err = v.RunSync(ctx, h.syncValidators)
	if err != nil {
		return err
	}

	err = v.RunAsync(ctx, h.asyncValidators)
	if err != nil {
		return err
	}

	_, err = RunActions(ctx, h.afterActions)

	return err
}

func (h *HookHandler) IsConfigured(config *c.HooksConfig) bool {
	return true
}

func RunActions(ctx i.SyncContext, actions []Action) (bool, error) {
	for _, action := range actions {
		next, err := action(ctx)
		if err != nil || !next {
			return false, err
		}
	}

	return true, nil
}
