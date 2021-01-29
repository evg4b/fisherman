package hookfactory

import (
	"fisherman/internal/handling"
	"fisherman/internal/validation"
)

var NoAsyncValidators = []validation.AsyncValidator{}
var NoAfterActions = []handling.Action{}
