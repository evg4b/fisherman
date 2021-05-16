package expression

import (
	"fisherman/utils"

	"github.com/antonmedv/expr"
	"github.com/imdario/mergo"
)

type Engine interface {
	Eval(expression string, vars map[string]interface{}) (bool, error)
	EvalMap(expression string, variables map[string]interface{}) (map[string]interface{}, error)
}

type GoExpressionEngine struct {
	functions map[string]interface{}
}

func NewGoExpressionEngine() *GoExpressionEngine {
	return &GoExpressionEngine{
		// TODO: Add functions:
		// - filesChanged(...glob) bool
		// - filesExist(...glob) bool
		// - env(name string) string
		// - filesChangedRelativeTo(...glob, branch) bool
		functions: map[string]interface{}{
			"IsEmpty": utils.IsEmpty,
			"Extract": extract,
		},
	}
}

func (engine *GoExpressionEngine) Eval(expressionString string, vars map[string]interface{}) (bool, error) {
	env := map[string]interface{}{}
	err := mergo.Merge(&env, engine.functions)
	if err != nil {
		return false, err
	}

	err = mergo.Merge(&env, vars)
	if err != nil {
		return false, err
	}

	program, err := expr.Compile(
		expressionString,
		expr.Env(env),
		expr.AllowUndefinedVariables(),
		expr.AsBool())

	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}

func (engine *GoExpressionEngine) EvalMap(expressionString string, vars map[string]interface{}) (map[string]interface{}, error) {
	env := map[string]interface{}{}
	err := mergo.Merge(&env, engine.functions)
	if err != nil {
		return nil, err
	}

	err = mergo.Merge(&env, vars)
	if err != nil {
		return nil, err
	}

	program, err := expr.Compile(expressionString, expr.Env(env), expr.AllowUndefinedVariables())

	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return nil, err
	}

	return output.(map[string]interface{}), nil
}
