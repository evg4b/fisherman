package hookfactory

import (
	"fisherman/internal/handling"
	"fisherman/internal/validation"
)

var NoSyncValidators = []validation.SyncValidator{}
var NoAsyncValidators = []validation.AsyncValidator{}
var NoBeforeActions = []handling.Action{}
var NoAfterActions = []handling.Action{}
