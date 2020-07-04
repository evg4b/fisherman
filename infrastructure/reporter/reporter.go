package reporter

type Reporter interface {
	ValidationError(rule, message string)
	Error(message string)
	Info(message string)
	Debug(message string)
}
