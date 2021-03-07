package initialize

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildHook(t *testing.T) {
	testData := []struct {
		hookName     string
		binaryPath   string
		absolute     bool
		expectedHook string
	}{
		{
			hookName:   "applypatch-msg",
			binaryPath: "fisherman",
			expectedHook: strings.Join([]string{
				"#!/bin/sh",
				"# This is fisherman hook handler. Please DO NOT touch this file.",
				"fisherman handle --hook applypatch-msg $@",
			}, LineBreak),
			absolute: false,
		},
		{
			hookName:   "pre-commit",
			binaryPath: "/bin/usr/fisherman",
			expectedHook: strings.Join([]string{
				"#!/bin/sh",
				"# This is fisherman hook handler. Please DO NOT touch this file.",
				"'/bin/usr/fisherman' handle --hook pre-commit $@",
			}, LineBreak),
			absolute: true,
		},
		{
			hookName:   "pre-push",
			binaryPath: "C:\\bin\\usr\\fisherman.exe",
			expectedHook: strings.Join([]string{
				"#!/bin/sh",
				"# This is fisherman hook handler. Please DO NOT touch this file.",
				"'C:\\\\bin\\\\usr\\\\fisherman.exe' handle --hook pre-push $@",
			}, LineBreak),
			absolute: true,
		},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Build %s hook", tt.hookName), func(t *testing.T) {
			result := buildHook(tt.hookName, tt.binaryPath, tt.absolute)
			assert.Equal(t, result, []byte(tt.expectedHook))
		})
	}
}
