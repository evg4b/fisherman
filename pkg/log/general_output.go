//go:build !test
// +build !test

package log

import (
	"io"
	"os"
)

var generalOutput io.Writer = os.Stdout
