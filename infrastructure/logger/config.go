package logger

// OutputConfig is structure to configure logger
type OutputConfig struct {
	LogLevel LogLevel `yaml:"level"`
	Colors   bool     `yaml:"colors"`
}

// DefaultOutputConfig is default values for configuration
var DefaultOutputConfig = OutputConfig{
	LogLevel: Info,
	Colors:   true,
}
