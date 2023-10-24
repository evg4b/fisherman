package handling

import (
  "github.com/evg4b/fisherman/internal/configuration"
  "github.com/evg4b/fisherman/internal/rules"
)

type Rule = configuration.Rule

func getPreScriptRules(ruleCollection []Rule) []Rule {
	return filterRules(ruleCollection, func(r Rule) bool {
		return r.GetPosition() == rules.PreScripts
	})
}

func getPostScriptRules(ruleCollection []Rule) []Rule {
	return filterRules(ruleCollection, func(r Rule) bool {
		return r.GetPosition() == rules.PostScripts
	})
}

func getScriptRules(ruleCollection []Rule) []Rule {
	return filterRules(ruleCollection, func(r Rule) bool {
		return r.GetPosition() == rules.Scripts
	})
}

func filterRules(rules []Rule, predicate func(Rule) bool) []Rule {
	filteredRules := []Rule{}

	for _, rule := range rules {
		if predicate(rule) {
			filteredRules = append(filteredRules, rule)
		}
	}

	return filteredRules
}
