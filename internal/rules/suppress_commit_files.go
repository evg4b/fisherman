package rules

import (
	"context"
	"fisherman/internal/utils"
	"io"
	"path/filepath"

	"github.com/go-errors/errors"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-multierror"
)

const SuppressCommitFilesType = "suppress-commit-files"

type SuppressCommitFiles struct {
	BaseRule        `yaml:",inline"`
	Globs           []string `yaml:"globs"`
	RemoveFromIndex bool     `yaml:"remove-from-index"`
}

// nolint: cyclop
func (rule *SuppressCommitFiles) Check(_ context.Context, _ io.Writer) error {
	if len(rule.Globs) == 0 {
		return nil
	}

	repo := rule.BaseRule.repo
	files, err := repo.GetFilesInIndex()
	if err != nil {
		return err
	}

	matchedFiles := []string{}
	for _, glob := range rule.Globs {
		for _, file := range files {
			matched, err := filepath.Match(glob, file)
			if err != nil {
				return err
			}
			if matched {
				matchedFiles = append(matchedFiles, file)
			}
		}
	}

	var multiError *multierror.Error
	for _, file := range matchedFiles {
		if rule.RemoveFromIndex {
			err := repo.RemoveGlob(file)
			if err != nil && !errors.Is(err, git.ErrGlobNoMatches) {
				return err
			}
		} else {
			multiError = multierror.Append(multiError, rule.errorf("file %s can not be committed", file))
		}
	}

	return multiError.ErrorOrNil()
}

func (rule *SuppressCommitFiles) GetPosition() byte {
	return PostScripts
}

func (rule *SuppressCommitFiles) Compile(variables map[string]any) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplatesArray(rule.Globs, variables)
}
