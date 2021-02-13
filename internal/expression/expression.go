package expression

import (
	"github.com/Knetic/govaluate"
	"github.com/imdario/mergo"
)

type Engine interface {
	Eval(expression string) (bool, error)
	EvalMap(expression string, variables map[string]interface{}) (map[string]interface{}, error)
}

type GovaluateEngine struct {
	globalFunctions map[string]govaluate.ExpressionFunction
	globalVariables map[string]interface{}
}

func NewExpressionEngine(variables map[string]interface{}) *GovaluateEngine {
	return &GovaluateEngine{
		globalVariables: variables,
		// TODO: Add functions:
		// - filesChanged(...glob) bool
		// - filesExist(...glob) bool
		// - env(name string) string
		// - filesChangedRelativeTo(...glob, branch) bool
		globalFunctions: map[string]govaluate.ExpressionFunction{
			"IsEmpty": isEmpty,
		},
	}
}

func (engine *GovaluateEngine) Eval(expressionString string) (bool, error) {
	// TODO: add global and local function. This case need to configure unique functions for each hook
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expressionString, engine.globalFunctions)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(engine.globalVariables)
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

	combinedVars := map[string]interface{}{}
	err = mergo.MergeWithOverwrite(&combinedVars, engine.globalVariables)
	if err != nil {
		return nil, err
	}

	err = mergo.MergeWithOverwrite(&combinedVars, vars)
	if err != nil {
		return nil, err
	}

	result, err := expression.Evaluate(combinedVars)
	if err != nil {
		return nil, err
	}

	return result.(map[string]interface{}), nil
}
