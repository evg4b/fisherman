package expression

import (
	"errors"
	"fisherman/utils"
)

func isEmpty(arguments ...interface{}) (interface{}, error) {
	// TODO: find validation mechanism
	if len(arguments) != 1 {
		return nil, errors.New("incorrect arguments for isEmpty")
	}

	return utils.IsEmpty(arguments[0].(string)), nil
}
