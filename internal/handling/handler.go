package handling

import (
	"context"
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
	"fisherman/internal/expression"
	"fisherman/internal/rules"
	"fisherman/internal/utils"
	"io"
	"runtime"

	"github.com/go-errors/errors"
	"github.com/go-git/go-billy/v5"
	"github.com/hashicorp/go-multierror"
	"github.com/imdario/mergo"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

type CompilableConfig interface {
	Compile(engine expression.Engine, global map[string]any) (map[string]any, error)
}

type Handler interface {
	Handle(ctx context.Context) error
}

type HookHandler struct {
	engine       expression.Engine
	configs      *configuration.HooksConfig
	globalVars   map[string]any
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
		globalVars:      map[string]any{},
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

		rule.Configure(h.ruleOptions()...)
	}

	err = multiError.ErrorOrNil()
	if err != nil {
		return nil, errors.Errorf("%s hook: %v", hook, err)
	}

	h.rules = getPreScriptRules(config.Rules)
	h.scripts = getScriptRules(config.Rules)
	h.postScriptRules = getPostScriptRules(config.Rules)

	vars, err := h.getPredefinedVariables()
	if err != nil {
		return nil, err
	}

	err = mergo.MergeWithOverwrite(&vars, h.globalVars)
	if err != nil {
		return nil, err
	}

	h.globalVars, err = config.Compile(vars)
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

func (h *HookHandler) ruleOptions() []rules.RuleOption {
	return []rules.RuleOption{
		rules.WithCwd(h.cwd),
		rules.WithFileSystem(h.fs),
		rules.WithRepository(h.repo),
		rules.WithArgs(h.args),
		rules.WithEnv(h.env),
	}
}

func (h *HookHandler) getPredefinedVariables() (map[string]any, error) {
	gitUser, err := h.repo.GetUser()
	if err != nil {
		return nil, err
	}

	branch, err := h.repo.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	tag, err := h.repo.GetLastTag()
	if err != nil {
		return nil, err
	}

	return map[string]any{
		constants.UserEmailVariable:        gitUser.Email,
		constants.UserNameVariable:         gitUser.UserName,
		constants.FishermanVersionVariable: constants.Version,
		constants.CwdVariable:              h.cwd,
		constants.BranchNameVariable:       branch,
		constants.TagVariable:              tag,
		constants.OsVariable:               runtime.GOOS,
	}, nil
}
