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

func init() {
	var err error
	FileLog, err = NewFileLogger("app.log")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	FileLog.Info("Log File Created")
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
func (l *FileLogger) logMessage(level string, info string, messages ...interface{}) {
	levels := "[" + level + "] "
	message := fmt.Sprintf(info, messages...)
	log.SetOutput(l.LogFile)
	log.Println(levels, message)
}

// Info logs an info message to the log file with a timestamp and date.
func (l *FileLogger) Info(info string, messages ...interface{}) {
	l.logMessage("INFO", info, messages...)
}

// Warning logs a warning message to the log file with a timestamp and date.
func (l *FileLogger) Warning(info string, messages ...interface{}) {
	l.logMessage("WARNING", info, messages...)
}

// Error logs an error message to the log file with a timestamp and date.
func (l *FileLogger) Error(info string, messages ...interface{}) {
	l.logMessage("ERROR", info, messages...)
}

func LogClose() {
	FileLog.LogFile.Close()
}
