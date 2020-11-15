package hooks

import (
	"fisherman/internal/handling"
	"fmt"
	"strings"
)

type HandlerRegistration struct {
	Registered bool
	Handler    handling.Handler
}

var NotRegistered HandlerRegistration = HandlerRegistration{Registered: false}

var NotSupported HandlerRegistration = HandlerRegistration{
	Registered: true,
	Handler:    new(handling.NotSupportedHandler),
}

type HandlerList map[string]HandlerRegistration

func (handlers HandlerList) Get(name string) (handling.Handler, error) {
	if hook, ok := handlers[strings.ToLower(name)]; ok {
		if hook.Registered {
			return hook.Handler, nil
		}

		return nil, nil
	}

	return nil, fmt.Errorf("'%s' is not valid hook name", name)
}
