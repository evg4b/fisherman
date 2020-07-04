package rules

import (
	"fisherman/common/validation"
	"fisherman/common/vcs"
	"fisherman/infrastructure/reporter"
	"fmt"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"sync"
)

type CommitMsgHookConfig struct {
	NotEmpty     bool   `yaml:"not-empty,omitempty"`
	CommitRegexp string `yaml:"commit-regexp,omitempty"`
	CommitPrefix string `yaml:"commit-prefix,omitempty"`
	CommitSuffix string `yaml:"commit-suffix,omitempty"`
}

func (c *CommitMsgHookConfig) BuildRules() validation.ValidationRule {

	return func(wg *sync.WaitGroup, snapshot vcs.Snapshot, reporter reporter.Reporter) {
		defer wg.Done()
		info := snapshot.GetHookData()
		message := info["message"]

		validationRules := make([]validate.Validator, 10)
		if c.NotEmpty == false {
			validationRules = append(validationRules, &validators.StringIsPresent{
				Name:    "Commit message",
				Field:   message,
				Message: "Commit message should not be empty.",
			})
		}

		validationRules = append(validationRules, &validators.StringIsPresent{
			Name:    "Commit prefix",
			Field:   message,
			Message: fmt.Sprintf("Commit message should have prefix '%s'."),
		})

		// errors := validate.Validate(validationRules)
		errors := validate.Validate(
			&validators.StringIsPresent{Field: info[""], Name: "Name"},
			&validators.StringIsPresent{Field: info[""], Name: "Name"},
		)

		for key, errs := range errors.Errors {
			for _, errorMessage := range errs {
				reporter.ValidationError(key, errorMessage)
			}
		}
	}
}
