package utils_test

import (
	"bytes"
	"fisherman/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintGraphics(t *testing.T) {
	tests := []struct {
		name    string
		content string
		data    interface{}
		wantWr  string
	}{
		{name: "Print template without data", content: "Template", data: nil, wantWr: "Template"},
		{name: "Print template with empty data map", content: "Template", data: make(map[string]string), wantWr: "Template"},
		{
			name:    "Print template with correct data",
			content: "Template [{{.Demo}}] = {{.Test}}",
			data: map[string]string{
				"Demo": "this is demo",
				"Test": "this is test",
			},
			wantWr: "Template [this is demo] = this is test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := &bytes.Buffer{}
			utils.PrintGraphics(wr, tt.content, tt.data)
			assert.Equal(t, tt.wantWr, wr.String())
		})
	}
}
