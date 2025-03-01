package logger

import (
	"log"
	"os"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type StandardLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

// NewStandardLogger initializes a new logger with different log levels.
func NewStandardLogger() *StandardLogger {
	// Create a log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return &StandardLogger{
		debugLogger: log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Implement the Logger interface methods
func (l *StandardLogger) Debug(msg string) {
	l.debugLogger.Println(msg)
}

func (l *StandardLogger) Info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *StandardLogger) Warn(msg string) {
	l.warnLogger.Println(msg)
}

func (l *StandardLogger) Error(msg string) {
	l.errorLogger.Println(msg)
}
