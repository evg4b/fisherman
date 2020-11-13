package initialize

import (
	"fisherman/constants"
	"fmt"
	"strings"
)

const LineBreak = "\n"

func buildHook(binaryPath, hookName string) string {
	return rows([]string{
		"#!/bin/sh",
		fmt.Sprintf("# This is %s hook handler. Please DO NOT touch this file.", constants.AppName),
		command([]string{escape(binaryPath), "handle", "--hook", hookName, "$@"}),
	})
}

func escape(path string) string {
	return strings.ReplaceAll(path, "\\", "\\\\")
}

func command(params []string) string {
	return strings.Join(params, " ")
}

func rows(rows []string) string {
	return strings.Join(rows, LineBreak)
}
