package expression

import (
	"fisherman/internal/constants"
	"fisherman/internal/utils"

	"github.com/antonmedv/expr"
	"github.com/imdario/mergo"
)

type Engine interface {
	Eval(expression string, vars map[string]interface{}) (bool, error)
}

type GoExpressionEngine struct{}

func NewGoExpressionEngine() *GoExpressionEngine {
	return &GoExpressionEngine{}
}

func (*GoExpressionEngine) Eval(expString string, vars map[string]interface{}) (bool, error) {
	env := map[string]interface{}{}

	if err := mergo.Merge(&env, resolveFunctions(vars)); err != nil {
		return false, err
	}

	if err := mergo.Merge(&env, vars); err != nil {
		return false, err
	}

	expression, err := expr.Compile(expString, expr.Env(env), expr.AllowUndefinedVariables(), expr.AsBool())
	if err != nil {
		return false, err
	}

	output, err := expr.Run(expression, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}

func resolveFunctions(vars map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"IsEmpty":   utils.IsEmpty,
		"IsWindows": func() bool { return vars[constants.OsVariable] == "windows" },
		"IsLinux":   func() bool { return vars[constants.OsVariable] == "linux" },
		"IsMac":     func() bool { return vars[constants.OsVariable] == "darwin" },
	}
}
