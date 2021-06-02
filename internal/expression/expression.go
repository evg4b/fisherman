package expression

import (
	"fisherman/internal/utils"

	"github.com/antonmedv/expr"
	"github.com/imdario/mergo"
)

type Engine interface {
	Eval(expression string, vars map[string]interface{}) (bool, error)
}

type GoExpressionEngine struct {
	functions map[string]interface{}
}

func NewGoExpressionEngine() *GoExpressionEngine {
	return &GoExpressionEngine{
		functions: map[string]interface{}{
			"IsEmpty": utils.IsEmpty,
		},
	}
}

func (engine *GoExpressionEngine) Eval(expString string, vars map[string]interface{}) (bool, error) {
	env := map[string]interface{}{}

	err := mergo.Merge(&env, engine.functions)
	if err != nil {
		return false, err
	}

	err = mergo.Merge(&env, vars)
	if err != nil {
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
