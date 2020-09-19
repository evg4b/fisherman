package init

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildHook(t *testing.T) {
	testData := []struct {
		hookName     string
		expectedHook string
	}{
		{"applypatch-msg", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook applypatch-msg $@"},
		{"commit-msg", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook commit-msg $@"},
		{"fsmonitor-watchman", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook fsmonitor-watchman $@"},
		{"post-update", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook post-update $@"},
		{"pre-applypatch", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook pre-applypatch $@"},
		{"pre-commit", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook pre-commit $@"},
		{"pre-push", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook pre-push $@"},
		{"pre-rebase", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook pre-rebase $@"},
		{"pre-receive", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook pre-receive $@"},
		{"prepare-commit-msg", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook prepare-commit-msg $@"},
		{"update", "#!/bin/sh\n# This is fisherman hook handler. Please DO NOT touch this file.\nfisherman handle --hook update $@"},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Build %s hook", tt.hookName), func(t *testing.T) {
			result := buildHook("fisherman", tt.hookName)
			assert.Equal(t, result, tt.expectedHook)
		})
	}
}
