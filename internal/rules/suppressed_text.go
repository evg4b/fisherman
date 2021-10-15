package rules

import (
	"fisherman/internal"
	"fisherman/internal/utils"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/go-multierror"
)

const SuppressedTextType = "suppressed-text"

type SuppressedText struct {
	BaseRule      `yaml:",inline"`
	Substrings    []string `yaml:"substrings"`
	ExcludedGlobs []string `yaml:"exclude"`
}

func (rule *SuppressedText) Check(ctx internal.ExecutionContext, _ io.Writer) error {
	repo := ctx.Repository()
	changes, err := repo.GetIndexChanges()
	if err != nil {
		return err
	}

	var validationErrors multierror.Error
	validationErrors.ErrorFormat = plainErrorFormatter

	for file, fileChanges := range changes {
		matched, err := utils.MatchToGlobs(rule.ExcludedGlobs, file)
		if err != nil {
			return err
		}

		if matched {
			continue
		}

		addedLines := fileChanges.Added()
		for _, suppressed := range rule.Substrings {
			for _, lineChange := range addedLines {
				if strings.Contains(lineChange.Change, suppressed) {
					err := fmt.Errorf("file '%s' should not contains '%s'", file, suppressed)
					validationErrors = *multierror.Append(&validationErrors, err)

					break
				}
			}
		}
	}

	return validationErrors.ErrorOrNil()
}

func (rule *SuppressedText) GetPosition() byte {
	return PostScripts
}

func (rule *SuppressedText) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplatesArray(rule.Substrings, variables)
	utils.FillTemplatesArray(rule.ExcludedGlobs, variables)
}
