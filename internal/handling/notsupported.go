package handling

import (
	"fisherman/config"
	"fisherman/constants"
	"fmt"
)

type NotSupportedHandler struct{}

func (*NotSupportedHandler) Handle(args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

func (*NotSupportedHandler) IsConfigured(*config.HooksConfig) bool {
	return true
}
