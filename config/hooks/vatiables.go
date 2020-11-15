package hooks

import (
	"fisherman/utils"
	"fmt"
	"regexp"
)

type VariablesExtractor interface {
	GetFromBranch(branchName string) (map[string]interface{}, error)
	GetFromTag(tag string) (map[string]interface{}, error)
}

type Variables struct {
	FromBranch  string `yaml:"from-branch,omitempty"`
	FromLastTag string `yaml:"from-last-tag,omitempty"`
}

func (config *Variables) GetFromBranch(branchName string) (map[string]interface{}, error) {
	return ejectFromString(branchName, config.FromBranch)
}

func (config *Variables) GetFromTag(tag string) (map[string]interface{}, error) {
	return ejectFromString(tag, config.FromLastTag)
}

func ejectFromString(tag, expressionString string) (map[string]interface{}, error) {
	variables := make(map[string]interface{})

	if utils.IsNotEmpty(expressionString) && utils.IsNotEmpty(tag) {
		reg, err := regexp.Compile(expressionString)
		if err != nil {
			return nil, err
		}

		match := reg.FindStringSubmatch(tag)
		if match == nil {
			return nil, fmt.Errorf("filed match '%s' to expression '%s'", tag, expressionString)
		}

		for i, name := range reg.SubexpNames() {
			if utils.IsNotEmpty(name) {
				variables[name] = match[i]
			}
		}
	}

	return variables, nil
}
