package initialize

import (
	"fisherman/constants"
	"fmt"
	"strings"
)

func buildHook(binaryPath, hookName string) string {
	return rows([]string{
		"#!/bin/sh",
		fmt.Sprintf("# This is %s hook handler. Please DO NOT touch this file.", constants.AppName),
		command([]string{binaryPath, "handle", "--hook", hookName, "$@"}),
	})
}

func command(params []string) string {
	return strings.Join(params, " ")
}

func rows(rows []string) string {
	return strings.Join(rows, "\n")
}
