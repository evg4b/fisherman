package handling

import (
	"context"
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/expression"
	"fisherman/internal/rules"
	"fisherman/internal/utils"
	"io"

	"github.com/go-errors/errors"
	"github.com/go-git/go-billy/v5"
	"github.com/hashicorp/go-multierror"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

type CompilableConfig interface {
	Compile(engine expression.Engine, global map[string]interface{}) (map[string]interface{}, error)
}

type Handler interface {
	Handle(ctx context.Context) error
}

type HookHandler struct {
	engine       expression.Engine
	configs      *configuration.HooksConfig
	globalVars   map[string]interface{}
	cwd          string
	fs           billy.Filesystem
	repo         internal.Repository
	args         []string
	env          []string
	workersCount uint
	output       io.Writer

	rules           []configuration.Rule
	scripts         []configuration.Rule
	postScriptRules []configuration.Rule
}

func NewHookHandler(hook string, options ...handlerOptions) (*HookHandler, error) {
	h := &HookHandler{
		configs:         &configuration.HooksConfig{},
		globalVars:      map[string]interface{}{},
		args:            []string{},
		env:             []string{},
		workersCount:    1,
		rules:           []configuration.Rule{},
		scripts:         []configuration.Rule{},
		postScriptRules: []configuration.Rule{},
	}

	for _, option := range options {
		option(h)
	}

	config, err := getConfig(hook, h.configs)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, ErrNotPresented
	}

	var multiError *multierror.Error
	for _, rule := range config.Rules {
		if !utils.Contains(allowedHooks[hook], rule.GetType()) {
			multiError = multierror.Append(multiError, errors.Errorf("rule %s is not allowed", rule.GetType()))

			continue
		}

		rule.Init(
			rules.WithCwd(h.cwd),
			rules.WithFileSystem(h.fs),
			rules.WithRepository(h.repo),
			rules.WithArgs(h.args),
			rules.WithEnv(h.env),
		)
	}

	err = multiError.ErrorOrNil()
	if err != nil {
		return nil, errors.Errorf("%s hook: %v", hook, err)
	}

	h.rules = getPreScriptRules(config.Rules)
	h.scripts = getScriptRules(config.Rules)
	h.postScriptRules = getPostScriptRules(config.Rules)
	h.globalVars, err = config.Compile(h.globalVars)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *HookHandler) Handle(ctx context.Context) error {
	err := h.run(ctx, h.rules)
	if err != nil {
		return err
	}

	err = h.run(ctx, h.scripts)
	if err != nil {
		return err
	}

	return h.run(ctx, h.postScriptRules)
}

func (h *HookHandler) run(ctx context.Context, rules []configuration.Rule) error {
	filteredRules := []configuration.Rule{}
	for _, rule := range rules {
		shouldAdd := true

		condition := rule.GetContition()
		if !utils.IsEmpty(condition) {
			var err error
			shouldAdd, err = h.engine.Eval(condition, h.globalVars)
			if err != nil {
				return err
			}
		}

		if shouldAdd {
			filteredRules = append(filteredRules, rule)
		}
	}

	return h.runRules(ctx, filteredRules)
}
