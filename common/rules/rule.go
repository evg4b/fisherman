package rules

import "fisherman/common/validation"

type RuleBuilder interface {
	BuildRules() []validation.ValidationRule
}
