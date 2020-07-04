package init

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fisherman/infrastructure/io"
	"fmt"
	"path/filepath"
)

const hookDir = "hooks"
const gitDir = ".git"

func writeHooks(path string, appInfo *context.AppInfo, accessor io.FileAccessor, force bool) error {
	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(path, gitDir, hookDir, hookName)
		if !force && accessor.FileExist(hookPath) {
			return fmt.Errorf("file %s already exists", hookPath)
		}

		bin := "fisherman"
		if !appInfo.IsRegisteredInPath {
			bin = appInfo.AppPath
		}

		content := buildHook(bin, hookName)
		err := accessor.WriteFile(hookPath, content)
		if err != nil {
			return err
		}
	}
	return nil
}
