package init

import "strings"

func buildHook(binaryPath, hookName string) string {
	return rows([]string{
		"#!/bin/sh",
		command([]string{binaryPath, "handle", "--hook", hookName, "$@"}),
	})
}

func command(params []string) string {
	return strings.Join(params, " ")
}

func rows(rows []string) string {
	return strings.Join(rows, "\n")
}
