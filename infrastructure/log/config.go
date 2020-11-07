package log

import (
	"errors"
	"strings"
)

// OutputConfig is structure to configure logger
type OutputConfig struct {
	LogLevel Level `yaml:"level"`
	Colors   bool  `yaml:"colors"`
}

// DefaultOutputConfig is default values for configuration
var DefaultOutputConfig = OutputConfig{
	LogLevel: InfoLevel,
	Colors:   true,
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (c *OutputConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var config struct {
		LogLevel string `yaml:"level"`
		Colors   bool   `yaml:"colors"`
	}

	if err := unmarshal(&config); err != nil {
		return err
	}

	level, err := parselogLevel(config.LogLevel)
	if err != nil {
		return err
	}

	c.Colors = config.Colors
	c.LogLevel = level

	return nil
}

func parselogLevel(level string) (Level, error) {
	if strings.EqualFold(level, "error") {
		return ErrorLevel, nil
	}

	if strings.EqualFold(level, "debug") {
		return DebugLevel, nil
	}

	if strings.EqualFold(level, "info") {
		return InfoLevel, nil
	}

	if strings.EqualFold(level, "none") {
		return NoneLevel, nil
	}

	return NoneLevel, errors.New("incorrect log level")
}
