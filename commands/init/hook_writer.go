package init

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fisherman/infrastructure/io"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

const hookDir = "hooks"
const gitDir = ".git"

func writeHooks(appInfo *context.AppInfo, accessor io.FileAccessor, force bool) error {
	var result *multierror.Error

	if !force {
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(appInfo.Cwd, gitDir, hookDir, hookName)
			if accessor.FileExist(hookPath) {
				multierror.Append(result, fmt.Errorf("File %s already exists", hookPath))
			}
		}
	}

	if result != nil {
		return result
	}

	bin := constants.AppName
	if !appInfo.IsRegisteredInPath {
		bin = appInfo.AppPath
	}

	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(appInfo.Cwd, gitDir, hookDir, hookName)
		content := buildHook(bin, hookName)
		err := accessor.WriteFile(hookPath, content)
		if err != nil {
			return err
		}
	}

	return nil
}
