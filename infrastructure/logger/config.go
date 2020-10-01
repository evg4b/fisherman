package logger

import (
	"errors"
	"strings"
)

// OutputConfig is structure to configure logger
type OutputConfig struct {
	LogLevel LogLevel `yaml:"level"`
	Colors   bool     `yaml:"colors"`
}

// DefaultOutputConfig is default values for configuration
var DefaultOutputConfig = OutputConfig{
	LogLevel: InfoLevel,
	Colors:   true,
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (config *OutputConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var outputConfig struct {
		LogLevel string `yaml:"level"`
		Colors   bool   `yaml:"colors"`
	}

	if err := unmarshal(&outputConfig); err != nil {
		return err
	}

	level, err := parselogLevel(outputConfig.LogLevel)
	if err != nil {
		return err
	}

	config.Colors = outputConfig.Colors
	config.LogLevel = level

	return nil
}

func parselogLevel(level string) (LogLevel, error) {
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
