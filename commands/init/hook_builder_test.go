package init

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildHook(t *testing.T) {
	testData := []struct {
		hookName     string
		expectedHook string
	}{
		{"applypatch-msg", "#!/bin/sh\nfisherman handle --hook applypatch-msg $@"},
		{"commit-msg", "#!/bin/sh\nfisherman handle --hook commit-msg $@"},
		{"fsmonitor-watchman", "#!/bin/sh\nfisherman handle --hook fsmonitor-watchman $@"},
		{"post-update", "#!/bin/sh\nfisherman handle --hook post-update $@"},
		{"pre-applypatch", "#!/bin/sh\nfisherman handle --hook pre-applypatch $@"},
		{"pre-commit", "#!/bin/sh\nfisherman handle --hook pre-commit $@"},
		{"pre-push", "#!/bin/sh\nfisherman handle --hook pre-push $@"},
		{"pre-rebase", "#!/bin/sh\nfisherman handle --hook pre-rebase $@"},
		{"pre-receive", "#!/bin/sh\nfisherman handle --hook pre-receive $@"},
		{"prepare-commit-msg", "#!/bin/sh\nfisherman handle --hook prepare-commit-msg $@"},
		{"update", "#!/bin/sh\nfisherman handle --hook update $@"},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Build %s hook", tt.hookName), func(t *testing.T) {
			result := buildHook("fisherman", tt.hookName)
			assert.Equal(t, result, tt.expectedHook)
		})
	}
}
