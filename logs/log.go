package logs

import (
	"fmt"
	"log"
	"os"
)

var FileLog *FileLogger

type FileLogger struct {
	LogFile *os.File
}

// NewFileLogger creates a new FileLogger instance and opens the log file.
func NewFileLogger(logFilePath string) (*FileLogger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileLogger{logFile}, nil
}

// logMessage logs a message to the log file with the specified log level.
func (l *FileLogger) logMessage(level string, messages ...interface{}) {
	message := "[" + level + "] "
	message += fmt.Sprint(messages...)
	log.SetOutput(l.LogFile)
	log.Println(message)
}

// Info logs an info message to the log file with a timestamp and date.
func (l *FileLogger) Info(messages ...interface{}) {
	l.logMessage("INFO", messages...)
}

// Warning logs a warning message to the log file with a timestamp and date.
func (l *FileLogger) Warning(messages ...interface{}) {
	l.logMessage("WARNING", messages...)
}

// Error logs an error message to the log file with a timestamp and date.
func (l *FileLogger) Error(messages ...interface{}) {
	l.logMessage("ERROR", messages...)
}

func LogInstance() {
	var err error
	FileLog, err = NewFileLogger("app.log")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	FileLog, err = NewFileLogger("app.log")
	FileLog.Info("Log File Created")
}
func LogClose() {
	FileLog.LogFile.Close()
}
