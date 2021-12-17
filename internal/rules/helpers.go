package rules

import (
	"fisherman/internal/utils"
	"fmt"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

func plainErrorFormatter(e []error) string {
	parts := []string{}
	for _, err := range e {
		parts = append(parts, err.Error())
	}

	return strings.Join(parts, "\n")
}

func getEncoding(encodingName string) (encoding.Encoding, error) {
	if utils.IsEmpty(encodingName) {
		return encoding.Nop, nil
	}

	if enc, err := ianaindex.IANA.Encoding(encodingName); err == nil {
		return enc, nil
	}

	return nil, fmt.Errorf("'%s' is unknown encoding", encodingName)
}
