package actions

import (
	"errors"
	"fisherman/internal"
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-multierror"
)

type SuppresCommitFilesSections struct {
	Globs           []string `yaml:"globs"`
	RemoveFromIndex bool     `yaml:"remove-from-index"`
}

func (section *SuppresCommitFilesSections) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var globs []string
	err := unmarshal(&globs)
	if err == nil {
		section.Globs = globs

		return nil
	}

	type plain SuppresCommitFilesSections

	return unmarshal((*plain)(section))
}

func SuppresCommitFiles(ctx internal.ExecutionContext, section SuppresCommitFilesSections) (bool, error) {
	if len(section.Globs) == 0 {
		return true, nil
	}

	repo := ctx.Repository()
	files, err := repo.GetFilesInIndex()
	if err != nil {
		return false, err
	}

	matchedFiles := []string{}
	for _, glob := range section.Globs {
		for _, file := range files {
			matched, err := filepath.Match(glob, file)
			if err != nil {
				return false, err
			}
			if matched {
				matchedFiles = append(matchedFiles, file)
			}
		}
	}

	var multiError *multierror.Error
	for _, file := range matchedFiles {
		if section.RemoveFromIndex {
			err := repo.RemoveGlob(file)
			if err != nil && !errors.Is(err, git.ErrGlobNoMatches) {
				return false, err
			}
		} else {
			multiError = multierror.Append(multiError, fmt.Errorf("file %s can not be committed", file))
		}
	}

	err = multiError.ErrorOrNil()
	if err != nil {
		return true, err
	}

	return true, nil
}
