package expression

import (
	"github.com/Knetic/govaluate"
)

type Expression interface {
	Eval(expression string) (bool, error)
}

type Engine struct {
	functions map[string]govaluate.ExpressionFunction
	variables map[string]interface{}
}

func NewExpressionEngine(variables map[string]interface{}) *Engine {
	return &Engine{
		variables: variables,
		// TODO: Add functions:
		// - filesChanged(...glob)
		// - filesExist(...glob)
		// - env(name string)
		functions: map[string]govaluate.ExpressionFunction{
			"IsEmpty":    isEmpty,
			"IsNotEmpty": isNotEmpty,
		},
	}
}

func (engine *Engine) Eval(expressionString string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expressionString, engine.functions)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(engine.variables)
	if err != nil {
		return false, err
	}

	// TODO: Add casting to bool (https://github.com/spf13/cast/blob/8d17101741c81653ee960aa20f9febb31f1218aa/caste.go#L74)
	return result.(bool), nil
}
