package init

import (
	"fisherman/constants"
	"fisherman/infrastructure/io"
	"fmt"
	"path/filepath"
)

const hookDir = "hooks"
const gitDir = ".git"

func WriteHooks(path string, accessor io.FileAccessor, force bool) error {
	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(path, gitDir, hookDir, hookName)
		if !force && accessor.FileExist(hookPath) {
			return fmt.Errorf("file %s alrady exests", hookPath)
		}
		content := buildHook("fisherman", hookName)
		err := accessor.WriteFile(hookPath, content)
		if err != nil {
			return err
		}
	}
	return nil
}
