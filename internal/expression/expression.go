package expression

import (
	"github.com/Knetic/govaluate"
	"github.com/imdario/mergo"
)

type Engine interface {
	Eval(expression string, vars map[string]interface{}) (bool, error)
	EvalMap(expression string, variables map[string]interface{}) (map[string]interface{}, error)
}

type GovaluateEngine struct {
	globalFunctions map[string]govaluate.ExpressionFunction
}

func NewExpressionEngine() *GovaluateEngine {
	return &GovaluateEngine{
		// TODO: Add functions:
		// - filesChanged(...glob) bool
		// - filesExist(...glob) bool
		// - env(name string) string
		// - filesChangedRelativeTo(...glob, branch) bool
		// TODO: [Next realise] provide Defined function
		globalFunctions: map[string]govaluate.ExpressionFunction{
			"IsEmpty": isEmpty,
		},
	}
}

func (engine *GovaluateEngine) Eval(expressionString string, vars map[string]interface{}) (bool, error) {
	// TODO: add global and local function. This case need to configure unique functions for each hook
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expressionString, engine.globalFunctions)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(vars)
	if err != nil {
		return false, err
	}

	// TODO: Add casting to bool (https://github.com/spf13/cast/blob/8d17101741c81653ee960aa20f9febb31f1218aa/caste.go#L74)
	return result.(bool), nil
}

func (engine *GovaluateEngine) EvalMap(expr string, vars map[string]interface{}) (map[string]interface{}, error) {
	functions := map[string]govaluate.ExpressionFunction{
		"Extract": extract,
	}

	err := mergo.Merge(&functions, engine.globalFunctions)
	if err != nil {
		return nil, err
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
	if err != nil {
		return nil, err
	}

	result, err := expression.Evaluate(vars)
	if err != nil {
		return nil, err
	}

	return result.(map[string]interface{}), nil
}
