package handling

import (
	c "fisherman/config"
	"fisherman/infrastructure/log"
	v "fisherman/internal/validation"
	"io"
)

type Handler interface {
	IsConfigured(config *c.HooksConfig) bool
	Handle(args []string) error
}

type Action = func(v.SyncValidationContext) (bool, error)

type HookHandler struct {
	contextFactory  func(args []string, output io.Writer) *v.ValidationContext
	beforeActions   []Action
	syncValidators  []v.SyncValidator
	asyncValidators []v.AsyncValidator
	afterActions    []Action
}

func NewHookHandler(
	contextFactory func(args []string, output io.Writer) *v.ValidationContext,
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

func RunActions(ctx v.SyncValidationContext, actions []Action) (bool, error) {
	for _, action := range actions {
		next, err := action(ctx)
		if err != nil || !next {
			return false, err
		}
	}

	return true, nil
}
