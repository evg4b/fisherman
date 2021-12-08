package handling

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/expression"
	"fisherman/internal/utils"
)

type Handler interface {
	Handle(ctx internal.ExecutionContext) error
}

type Action = func(internal.ExecutionContext) (bool, error)

type HookHandler struct {
	Engine          expression.Engine
	Rules           []configuration.Rule
	Scripts         []configuration.Rule
	PostScriptRules []configuration.Rule
	GlobalVars      Variables
	WorkersCount    int
}

func (h *HookHandler) Handle(ctx internal.ExecutionContext) error {
	filterRules, err := h.filterRules(h.Rules)
	if err != nil {
		return err
	}

	err = h.runRules(ctx, filterRules)
	if err != nil {
		return err
	}

	filterScripts, err := h.filterRules(h.Scripts)
	if err != nil {
		return err
	}

	err = h.runRules(ctx, filterScripts)
	if err != nil {
		return err
	}

	filterPostScriptRules, err := h.filterRules(h.PostScriptRules)
	if err != nil {
		return err
	}

	return h.runRules(ctx, filterPostScriptRules)
}

func (h *HookHandler) filterRules(rules []configuration.Rule) ([]configuration.Rule, error) {
	filteredRules := []configuration.Rule{}
	for _, rule := range rules {
		shouldAdd := true

		condition := rule.GetContition()
		if !utils.IsEmpty(condition) {
			var err error
			shouldAdd, err = h.Engine.Eval(condition, h.GlobalVars)
			if err != nil {
				return nil, err
			}
		}

		if shouldAdd {
			filteredRules = append(filteredRules, rule)
		}
	}

	return filteredRules, nil
}
