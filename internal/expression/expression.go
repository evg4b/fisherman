package expression

import (
	"fmt"
	"reflect"

	"github.com/Knetic/govaluate"
)

type Engine interface {
	Eval(expression string, vars map[string]interface{}) (bool, error)
	EvalMap(expression string, variables map[string]interface{}) (map[string]interface{}, error)
}

type GovaluateEngine struct {
	functions map[string]govaluate.ExpressionFunction
}

func NewExpressionEngine() *GovaluateEngine {
	return &GovaluateEngine{
		// TODO: Add functions:
		// - filesChanged(...glob) bool
		// - filesExist(...glob) bool
		// - env(name string) string
		// - filesChangedRelativeTo(...glob, branch) bool
		functions: map[string]govaluate.ExpressionFunction{
			"IsEmpty": isEmpty,
			"Extract": extract,
			"Defined": defined,
		},
	}
}

func (engine *GovaluateEngine) Eval(expressionString string, vars map[string]interface{}) (bool, error) {
	// TODO: add global and local function. This case need to configure unique functions for each hook
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expressionString, engine.functions)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(vars)
	if err != nil {
		return false, err
	}

	result = indirect(result)

	switch b := result.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case float32:
		return b != 0, nil
	case int:
		return b != 0, nil
	case string:
		return len(b) > 0, nil
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", b, b)
	}
}

func (engine *GovaluateEngine) EvalMap(expr string, vars map[string]interface{}) (map[string]interface{}, error) {
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expr, engine.functions)
	if err != nil {
		return nil, err
	}

	result, err := expression.Evaluate(vars)
	if err != nil {
		return nil, err
	}

	return result.(map[string]interface{}), nil
}

func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}

	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		return a
	}

	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	return v.Interface()
}
