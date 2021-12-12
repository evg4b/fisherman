package handling

import (
	"fisherman/internal/expression"

	"github.com/go-errors/errors"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

type CompilableConfig interface {
	Compile(engine expression.Engine, global Variables) (Variables, error)
}

type (
	Variables   = map[string]interface{}
	hookBuilder = func(globalVars Variables) (Handler, error)
	builders    = map[string]hookBuilder
)
