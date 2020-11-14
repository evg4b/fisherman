package actions

import v "fisherman/internal/validation"

func AddToIndex(ctx v.SyncValidationContext, globs []string) (bool, error) {
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
