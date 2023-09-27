package logs

// Logger is a custom interface for logging.
type Logger interface {
	Info(message string)
	Warning(message string)
	Error(message string)
}

// FileLogger is an implementation of the Logger interface that logs to a file.
type FileLogger struct {
	filename string
}
