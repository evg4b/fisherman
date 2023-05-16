package rules

import (
	"context"
	"fisherman/internal/utils"
	"io"

	"github.com/go-errors/errors"

	"github.com/go-git/go-git/v5"
)

const AddToIndexType = "add-to-index"

type Glob struct {
	Pattern    string `yaml:"glob"`
	IsRequired bool   `yaml:"required"`
}

type AddToIndex struct {
	BaseRule `yaml:",inline"`
	Globs    []Glob `yaml:"globs"`
}

func (rule *AddToIndex) GetPosition() byte {
	return PostScripts
}

func (rule *AddToIndex) Check(_ context.Context, _ io.Writer) error {
	if len(rule.Globs) < 1 {
		return nil
	}

	repo := rule.BaseRule.repo
	for _, glob := range rule.Globs {
		if err := repo.AddGlob(glob.Pattern); err != nil {
			if errors.Is(err, git.ErrGlobNoMatches) {
				if !glob.IsRequired {
					continue
				}

				return rule.errorf("can't add files matching pattern %s", glob.Pattern)
			}

			return errors.Errorf("failed to add files matching pattern '%s' to the index: %w", glob.Pattern, err)
		}
	}

	return nil
}

func (rule *AddToIndex) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	for index := range rule.Globs {
		utils.FillTemplate(&rule.Globs[index].Pattern, variables)
	}
}
