package expression

import (
	"fisherman/utils"
	"fmt"
	"regexp"
)

func extract(arguments ...interface{}) map[string]interface{} {
	const extractArgumentCount = 2
	if len(arguments) != extractArgumentCount {
		panic("incorrect arguments for Extract")
	}

	source := arguments[0].(string)
	expression := arguments[1].(string)
	variables := make(map[string]interface{})

	if !utils.IsEmpty(expression) && !utils.IsEmpty(source) {
		reg, err := regexp.Compile(expression)
		if err != nil {
			panic(err)
		}

		match := reg.FindStringSubmatch(source)
		if match == nil {
			panic(fmt.Errorf("filed match '%s' to expression '%s'", source, expression))
		}

		for i, name := range reg.SubexpNames() {
			if !utils.IsEmpty(name) {
				variables[name] = match[i]
			}
		}
	}

	return variables
}
