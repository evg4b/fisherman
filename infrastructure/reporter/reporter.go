package reporter

// Reporter is interface to communicate with system output
type Reporter interface {
	ValidationError(rule, message string)
	Error(message string)
	Info(message string)
	Debug(message string)
	PrintGraphics(content string, data interface{})
}
