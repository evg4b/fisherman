package configuration

import (
	"fisherman/utils"
	"fmt"
	"regexp"
)

type Variables = map[string]interface{}

type VariablesExtractor interface {
	GetFromBranch(branchName string) (Variables, error)
	GetFromTag(tag string) (Variables, error)
}

type VariablesConfig struct {
	FromBranch  string `yaml:"from-branch,omitempty"`
	FromLastTag string `yaml:"from-last-tag,omitempty"`
}

func (config *VariablesConfig) GetFromBranch(branch string) (Variables, error) {
	return ejectFromString(branch, config.FromBranch)
}

func (config *VariablesConfig) GetFromTag(tag string) (Variables, error) {
	return ejectFromString(tag, config.FromLastTag)
}

func ejectFromString(tag, expression string) (Variables, error) {
	variables := make(Variables)

	if !utils.IsEmpty(expression) && !utils.IsEmpty(tag) {
		reg, err := regexp.Compile(expression)
		if err != nil {
			return nil, err
		}

		match := reg.FindStringSubmatch(tag)
		if match == nil {
			return nil, fmt.Errorf("filed match '%s' to expression '%s'", tag, expression)
		}

		for i, name := range reg.SubexpNames() {
			if !utils.IsEmpty(name) {
				variables[name] = match[i]
			}
		}
	}

	return variables, nil
}
