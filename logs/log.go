package logs

import (
	"log"
)

// SimpleLogger is a concrete implementation of the Logger interface using the standard log package.
type SimpleLogger struct{}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Info(message string) {
	log.Printf("[INFO] %s\n", message)
}

func (l *SimpleLogger) Warning(message string) {
	log.Printf("[WARNING] %s\n", message)
}

func (l *SimpleLogger) Error(message string) {
	log.Printf("[ERROR] %s\n", message)
}
