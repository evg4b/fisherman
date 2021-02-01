package hookfactory

import (
	"fisherman/configuration"
	"fisherman/internal/rules"
)

type Rule = configuration.Rule

func getPreScripts(ruleCollection []Rule) []Rule {
	return filterRules(ruleCollection, func(r Rule) bool {
		return r.GetPosition() == rules.PreScripts
	})
}

func getPostScriptRules(ruleCollection []Rule) []Rule {
	return filterRules(ruleCollection, func(r Rule) bool {
		return r.GetPosition() == rules.AfterScripts
	})
}

func getScriptRules(ruleCollection []Rule) []Rule {
	return filterRules(ruleCollection, func(r Rule) bool {
		return r.GetPosition() == rules.Scripts
	})
}

func filterRules(rules []Rule, predicate func(Rule) bool) []Rule {
	var filteredRules []Rule = []Rule{}

	for _, rule := range rules {
		if predicate(rule) {
			filteredRules = append(filteredRules, rule)
		}
	}

	return filteredRules
}
