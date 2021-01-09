package log

import (
	"errors"
	"strings"
)

type OutputConfig struct {
	LogLevel Level
	Colors   bool
}

var DefaultOutputConfig = OutputConfig{
	LogLevel: InfoLevel,
	Colors:   true,
}

var levelMatching = map[string]Level{
	"error": ErrorLevel,
	"debug": DebugLevel,
	"info":  InfoLevel,
	"none":  NoneLevel,
}

func (c *OutputConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	(*c) = DefaultOutputConfig

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
	value, ok := levelMatching[strings.Trim(level, " ")]
	if ok {
		return value, nil
	}

	return NoneLevel, errors.New("incorrect log level")
}
