package hooks

import (
	"fisherman/utils"
	"fmt"
	"regexp"
)

type Variables struct {
	FromBranch string `yaml:"from-branch,omitempty"`
}

func (c *Variables) GetFromBranch(branchName string) (map[string]interface{}, error) {
	variables := make(map[string]interface{})

	if utils.IsNotEmpty(c.FromBranch) {
		reg, err := regexp.Compile(c.FromBranch)
		if err != nil {
			return nil, err
		}

		match := reg.FindStringSubmatch(branchName)
		if match == nil {
			return nil, fmt.Errorf("filed match '%s' to expression '%s'", branchName, c.FromBranch)
		}

		for i, name := range reg.SubexpNames() {
			if utils.IsNotEmpty(name) {
				variables[name] = match[i]
			}
		}
	}

	return variables, nil
}
