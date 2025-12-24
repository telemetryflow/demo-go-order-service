// Package logger provides structured logging.
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents log level
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging
type Logger struct {
	level  Level
	output io.Writer
	fields map[string]interface{}
}

// New creates a new logger
func New(level string) *Logger {
	var l Level
	switch level {
	case "debug":
		l = DebugLevel
	case "info":
		l = InfoLevel
	case "warn":
		l = WarnLevel
	case "error":
		l = ErrorLevel
	default:
		l = InfoLevel
	}

	return &Logger{
		level:  l,
		output: os.Stdout,
		fields: make(map[string]interface{}),
	}
}

// WithField returns a new logger with an additional field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newFields := make(map[string]interface{}, len(l.fields)+1)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &Logger{
		level:  l.level,
		output: l.output,
		fields: newFields,
	}
}

// WithFields returns a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newFields := make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &Logger{
		level:  l.level,
		output: l.output,
		fields: newFields,
	}
}

func (l *Logger) log(level Level, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	entry := map[string]interface{}{
		"level":     level.String(),
		"message":   fmt.Sprintf(msg, args...),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	for k, v := range l.fields {
		entry[k] = v
	}

	data, _ := json.Marshal(entry)
	fmt.Fprintln(l.output, string(data))
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, msg, args...)
}

// Default logger instance
var defaultLogger = New("info")

// SetDefaultLogger sets the default logger
func SetDefaultLogger(l *Logger) {
	defaultLogger = l
}

// Debug logs using default logger
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

// Info logs using default logger
func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

// Warn logs using default logger
func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

// Error logs using default logger
func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}
