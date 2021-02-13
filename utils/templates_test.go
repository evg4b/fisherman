package utils_test

import (
	"bytes"
	"fisherman/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintGraphics(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		data           map[string]interface{}
		expectedOutput string
	}{
		{name: "Print template without data", content: "Template", data: nil, expectedOutput: "Template"},
		{
			name:           "Print template with empty data map",
			content:        "Template",
			data:           map[string]interface{}{},
			expectedOutput: "Template",
		},
		{
			name:    "Print template with correct data",
			content: "Template [{{Demo}}] = {{Test}}",
			data: map[string]interface{}{
				"Demo": "this is demo",
				"Test": "this is test",
			},
			expectedOutput: "Template [this is demo] = this is test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := &bytes.Buffer{}
			utils.PrintGraphics(wr, tt.content, tt.data)
			assert.Equal(t, tt.expectedOutput, wr.String())
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
		data     []map[string]interface{}
		expected string
	}{
		{
			name: "Should update template correctly",
			data: []map[string]interface{}{
				{"Test": "Test value"},
				{"Test2": "Test value2"},
			},
			template: "Template = {{Test}} + {{Test2}}",
			expected: "Template = Test value + Test value2",
		},
		{
			name: "Should skip template correctly",
			data: []map[string]interface{}{
				{"Test": "Test value"},
				{"Test2": "Test value2"},
			},
			template: "Template test",
			expected: "Template test",
		},
		{
			name: "Should skip template correctly",
			data: []map[string]interface{}{
				{"Test": "[value]"},
			},
			template: "{{Test}}={{Test}}={{Test}}",
			expected: "[value]=[value]=[value]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := utils.FillTemplate(tt.template, tt.data...)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
