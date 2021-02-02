package expression

import (
	"errors"
	"fisherman/utils"
	"fmt"
	"regexp"
)

func isEmpty(arguments ...interface{}) (interface{}, error) {
	// TODO: find validation mechanism
	const isEmptyArgumentCount = 1
	if len(arguments) != isEmptyArgumentCount {
		return nil, errors.New("incorrect arguments for IsEmpty")
	}

	return utils.IsEmpty(arguments[0].(string)), nil
}

func extract(arguments ...interface{}) (interface{}, error) {
	const extractArgumentCount = 2
	if len(arguments) != extractArgumentCount {
		return nil, errors.New("incorrect arguments for Extract")
	}

	source := arguments[0].(string)
	expression := arguments[1].(string)
	variables := make(map[string]interface{})

	if !utils.IsEmpty(expression) && !utils.IsEmpty(source) {
		reg, err := regexp.Compile(expression)
		if err != nil {
			return nil, err
		}

		match := reg.FindStringSubmatch(source)
		if match == nil {
			return nil, fmt.Errorf("filed match '%s' to expression '%s'", source, expression)
		}

		for i, name := range reg.SubexpNames() {
			if !utils.IsEmpty(name) {
				variables[name] = match[i]
			}
		}
	}

	return variables, nil
}
