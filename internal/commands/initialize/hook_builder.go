package initialize

import (
	"fmt"
	"github.com/evg4b/fisherman/internal/constants"
	"strings"
)

const LineBreak = "\n"

func buildHook(hookName, binaryPath string, absolute bool) []byte {
	template := "%s handle --hook %s $@"
	if absolute {
		template = "'%s' handle --hook %s $@"
	}

	data := strings.Join([]string{
		"#!/bin/sh",
		fmt.Sprintf("# This is %s hook handler. Please DO NOT touch this file.", constants.AppName),
		fmt.Sprintf(template, strings.ReplaceAll(binaryPath, "\\", "\\\\"), hookName),
	}, LineBreak)

	return []byte(data)
}
