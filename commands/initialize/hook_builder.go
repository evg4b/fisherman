package initialize

import (
	"fisherman/constants"
	"fmt"
	"strings"
)

const LineBreak = "\n"

func buildHook(hookName, binaryPath string, absolute bool) string {
	template := "%s handle --hook %s $@"
	if absolute {
		template = "'%s' handle --hook %s $@"
	}

	return strings.Join([]string{
		"#!/bin/sh",
		fmt.Sprintf("# This is %s hook handler. Please DO NOT touch this file.", constants.AppName),
		fmt.Sprintf(template, strings.ReplaceAll(binaryPath, "\\", "\\\\"), hookName),
	}, LineBreak)
}
