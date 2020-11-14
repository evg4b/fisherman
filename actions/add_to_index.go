package actions

import "fisherman/internal"

func AddToIndex(ctx internal.SyncContext, globs []string) (bool, error) {
	if len(globs) > 0 {
		repo := ctx.Repository()
		for _, glob := range globs {
			err := repo.AddGlob(glob)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}
