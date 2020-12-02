package hookfactory_test

import (
	"errors"
	"fisherman/internal/hookfactory"
	"fisherman/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerList_Get(t *testing.T) {
	handlers := hookfactory.HandlerList{
		"handler-1": hookfactory.NotRegistered,
		"handler-2": hookfactory.HandlerRegistration{
			Handler:    mocks.NewHandlerMock(t),
			Registered: true,
		},
	}

	tests := []struct {
		name        string
		handlerName string
		err         error
		hasHandler  bool
	}{
		{
			name:        "Not registered handler",
			handlerName: "handler-1",
			hasHandler:  false,
		},
		{
			name:        "Registered handler",
			handlerName: "handler-2",
			hasHandler:  true,
		},
		{
			name:        "Unknown handler",
			handlerName: "handler-3",
			hasHandler:  false,
			err:         errors.New("'handler-3' is not valid hook name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := handlers.Get(tt.handlerName)

			assert.Equal(t, tt.err, err)
			if tt.hasHandler {
				assert.IsType(t, &mocks.HandlerMock{}, handler)
			} else {
				assert.Nil(t, handler)
			}
		})
	}
}
