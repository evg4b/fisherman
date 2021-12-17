//go:build !windows
// +build !windows

package shell

import "golang.org/x/text/encoding"

func windowsEncoding() encoding.Encoding {
	return encoding.Nop
}
