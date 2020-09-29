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
		data    map[string]interface{}
		wantWr  string
	}{
		{name: "Print template without data", content: "Template", data: nil, wantWr: "Template"},
		{
			name:    "Print template with empty data map",
			content: "Template",
			data:    make(map[string]interface{}),
			wantWr:  "Template",
		},
		{
			name:    "Print template with correct data",
			content: "Template [{{Demo}}] = {{Test}}",
			data: map[string]interface{}{
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

func TestPrintGraphicsPanics(t *testing.T) {
	tests := []struct {
		name    string
		content string
		data    map[string]interface{}
	}{
		{name: "Panics when template is brocken", content: "Template{{Demo", data: nil},
		{name: "Panics when writer is nil", content: "Template", data: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				utils.PrintGraphics(nil, tt.content, tt.data)
			})
		})
	}
}

func TestFillTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expected string
	}{
		{name: "Should update template correctly", template: "Template = {{Test}}", expected: "Template = Test value"},
		{name: "Should skip template correctly", template: "Template", expected: "Template"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.FillTemplate(&tt.template, map[string]interface{}{
				"Test": "Test value",
			})

			assert.Equal(t, tt.expected, tt.template)
		})
	}
}
