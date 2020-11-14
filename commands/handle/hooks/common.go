package hooks

import (
	"io"

	h "fisherman/internal/handling"
	v "fisherman/internal/validation"
)

type ctxFactory = func(args []string, output io.Writer) *v.ValidationContext

var NoSyncValidators = []v.SyncValidator{}
var NoAsyncValidators = []v.AsyncValidator{}
var NoBeforeActions = []h.Action{}
var NoAfterActions = []h.Action{}
