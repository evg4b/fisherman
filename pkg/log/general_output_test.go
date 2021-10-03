//go:build test
// +build test

package log

import (
	"io"
)

var generalOutput io.Writer = io.Discard
