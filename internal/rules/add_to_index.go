package rules

import (
	"errors"
	"fisherman/internal"
	"io"

	"github.com/go-git/go-git/v5"
)

const AddToIndexType = "add-to-index"

type Glob struct {
	Glob       string `yaml:"glob"`
	IsRequired bool   `yaml:"required"`
}

type AddToIndex struct {
	BaseRule `mapstructure:",squash"`
	Globs    []Glob `mapstructure:"globs"`
}

func (rule *AddToIndex) GetPosition() byte {
	return PostScripts
}

func (config AddToIndex) Check(_ io.Writer, ctx internal.ExecutionContext) error {
	if len(config.Globs) > 0 {
		repo := ctx.Repository()
		for _, glob := range config.Globs {
			err := repo.AddGlob(glob.Glob)
			if err != nil {
				if errors.Is(err, git.ErrGlobNoMatches) && !glob.IsRequired {
					continue
				}

				return err
			}
		}
	}

	return nil
}
