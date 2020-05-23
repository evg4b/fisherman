package init

import "strings"

func buildHook(fishermanCommand, hookName string) string {
	return strings.Join([]string{fishermanCommand, "handle", "--hook", hookName}, " ")
}
