package log

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
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

func (c *OutputConfig) UnmarshalYAML(value *yaml.Node) error {
	(*c) = DefaultOutputConfig

	var config struct {
		LogLevel string `yaml:"level"`
		Colors   bool   `yaml:"colors"`
	}

	if err := value.Decode(&config); err != nil {
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
