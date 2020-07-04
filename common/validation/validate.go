package validation

import (
	"fisherman/common/vcs"
	"fisherman/infrastructure/reporter"
	"sync"
)

type ValidationResult struct {
	rule     string
	messages []string
}

type ValidationRule = func(wg *sync.WaitGroup, snapshot vcs.Snapshot, reporter reporter.Reporter)

func Validate(snapshot vcs.Snapshot, rules []ValidationRule) {
	var wg sync.WaitGroup
	for _, rule := range rules {
		wg.Add(1)
		go rule(&wg, snapshot, nil)
	}
	wg.Wait()
}
