package expression

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/imdario/mergo"
)

type Engine interface {
	Eval(expression string, vars map[string]interface{}) (bool, error)
}

type GoExpressionEngine struct {
	vm vm.VM
}

func NewGoExpressionEngine() *GoExpressionEngine {
	return &GoExpressionEngine{}
}

func (e *GoExpressionEngine) Eval(expString string, vars map[string]interface{}) (bool, error) {
	env := EnvVars{}
	var castedVars EnvVars = vars

	if err := mergo.Merge(&env, castedVars); err != nil {
		return false, err
	}

	expression, err := expr.Compile(expString, engineOptions(env)...)
	if err != nil {
		return false, err
	}

	output, err := e.vm.Run(expression, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}

func engineOptions(env EnvVars) []expr.Option {
	return []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(),
		expr.AsBool(),
		expr.Optimize(true),
	}
}
