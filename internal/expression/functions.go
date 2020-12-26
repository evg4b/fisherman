package expression

import (
	"errors"
	"fisherman/utils"
)

func isEmpty(arguments ...interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, errors.New("incorrect arguments for isEmpty")
	}

	return utils.IsEmpty(arguments[0].(string)), nil
}

func isNotEmpty(arguments ...interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, errors.New("incorrect arguments for isNotEmpty")
	}

	return utils.IsNotEmpty(arguments[0].(string)), nil
}
