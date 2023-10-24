package log_test

import (
	"testing"

	. "github.com/evg4b/fisherman/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestOutputConfig_UnmarshalYAML(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		tests := []struct {
			name   string
			config string
			colors bool
			level  Level
		}{
			{
				name:   "level is debug and colors is true",
				config: "level: debug\ncolors: true",
				colors: true,
				level:  DebugLevel,
			},
			{
				name:   "level is error and colors is true",
				config: "level: error\ncolors: true",
				colors: true,
				level:  ErrorLevel,
			},
			{
				name:   "level is none and colors is true",
				config: "level: none\ncolors: true",
				colors: true,
				level:  NoneLevel,
			},
			{
				name:   "level is info and colors is false",
				config: "level: info\ncolors: false",
				colors: false,
				level:  InfoLevel,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var config OutputConfig
				err := yaml.Unmarshal([]byte(tt.config), &config)

				require.NoError(t, err)
				assert.Equal(t, tt.colors, config.Colors)
				assert.Equal(t, tt.level, config.LogLevel)
			})
		}
	})

	t.Run("with errors", func(t *testing.T) {
		tests := []struct {
			name        string
			config      string
			expectedErr string
		}{
			{
				name:        "syntax error",
				config:      "level: debug\ncolor",
				expectedErr: "yaml: line 2: could not find expected ':'",
			},
			{
				name:        "incorrect log level constant",
				config:      "level: test",
				expectedErr: "incorrect log level",
			},
			{
				name:        "incrrect type",
				config:      "level: info\ncolors: 'test'",
				expectedErr: "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!str `test` into bool",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var config OutputConfig
				err := yaml.Unmarshal([]byte(tt.config), &config)
				require.EqualError(t, err, tt.expectedErr)
			})
		}
	})
}
