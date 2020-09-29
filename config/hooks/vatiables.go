package hooks

import (
	"fisherman/utils"
	"regexp"
)

type Variables struct {
	FromBranch string `yaml:"from-branch,omitempty"`
}

func (c *Variables) GetFromBranch(branchName string) (map[string]interface{}, error) {
	reg, err := regexp.Compile(c.FromBranch)
	if err != nil {
		return nil, err
	}

	match := reg.FindStringSubmatch(branchName)
	variables := make(map[string]interface{})
	subexpNames := reg.SubexpNames()
	for i, name := range subexpNames[1:] {
		if utils.IsNotEmpty(name) {
			variables[name] = match[i]
		}
	}

	return variables, err
}
