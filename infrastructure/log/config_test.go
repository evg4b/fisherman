package log_test

import (
	"fisherman/infrastructure/log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestOutputConfig_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name   string
		config string
		colors bool
		level  log.Level
	}{
		{name: "", config: "level: debug\ncolors: true", colors: true, level: log.DebugLevel},
		{name: "", config: "level: error\ncolors: true", colors: true, level: log.ErrorLevel},
		{name: "", config: "level: none\ncolors: true", colors: true, level: log.NoneLevel},
		{name: "", config: "level: info\ncolors: false", colors: false, level: log.InfoLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config log.OutputConfig
			err := yaml.Unmarshal([]byte(tt.config), &config)

			assert.NoError(t, err)
			assert.Equal(t, tt.colors, config.Colors)
			assert.Equal(t, tt.level, config.LogLevel)
		})
	}
}

func TestOutputConfig_UnmarshalYAMLFail(t *testing.T) {
	tests := []struct {
		name   string
		config string
		err    string
	}{
		{
			name:   "",
			config: "level: debug\ncolor",
			err:    "yaml: line 2: could not find expected ':'",
		},
		{
			name:   "",
			config: "level: test",
			err:    "incorrect log level",
		},
		{
			name:   "",
			config: "level: info\ncolors: 'test'",
			err:    "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!str `test` into bool",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config log.OutputConfig
			err := yaml.Unmarshal([]byte(tt.config), &config)
			assert.EqualError(t, err, tt.err)
		})
	}
}
