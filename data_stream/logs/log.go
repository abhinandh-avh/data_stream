package logs

// Logger is a custom interface for logging.
type Logger interface {
	Info(message string)
	Warning(message string)
	Error(message string)
}

