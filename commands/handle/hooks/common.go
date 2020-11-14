package hooks

import (
	h "fisherman/internal/handling"
	v "fisherman/internal/validation"
)

var NoSyncValidators = []v.SyncValidator{}
var NoAsyncValidators = []v.AsyncValidator{}
var NoBeforeActions = []h.Action{}
var NoAfterActions = []h.Action{}
