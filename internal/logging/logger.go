package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Level represents log level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

// Logger provides structured logging
type Logger struct {
	mu     sync.Mutex
	output io.Writer
	level  Level
	prefix string
	flags  int
}

// Global default logger
var defaultLogger *Logger
var once sync.Once

// Default returns the default logger instance
func Default() *Logger {
	once.Do(func() {
		defaultLogger = New(os.Stderr, LevelInfo, "[Laravel MCP] ", log.LstdFlags|log.Lshortfile)
	})
	return defaultLogger
}

// New creates a new logger
func New(output io.Writer, level Level, prefix string, flags int) *Logger {
	return &Logger{
		output: output,
		level:  level,
		prefix: prefix,
		flags:  flags,
	}
}

// SetLevel changes the logging level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// log writes a log message if level is enabled
func (l *Logger) log(level Level, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	levelStr := levelNames[level]
	message := fmt.Sprintf(format, v...)

	fmt.Fprintf(l.output, "%s%s [%s] %s\n", l.prefix, timestamp, levelStr, message)
}

// Debug logs at debug level
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(LevelDebug, format, v...)
}

// Info logs at info level
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(LevelInfo, format, v...)
}

// Warn logs at warn level
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(LevelWarn, format, v...)
}

// Error logs at error level
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(LevelError, format, v...)
}

// WithPrefix creates a new logger with additional prefix
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		output: l.output,
		level:  l.level,
		prefix: l.prefix + prefix + " ",
		flags:  l.flags,
	}
}

// Package-level convenience functions using default logger

// Debug logs at debug level using default logger
func Debug(format string, v ...interface{}) {
	Default().Debug(format, v...)
}

// Info logs at info level using default logger
func Info(format string, v ...interface{}) {
	Default().Info(format, v...)
}

// Warn logs at warn level using default logger
func Warn(format string, v ...interface{}) {
	Default().Warn(format, v...)
}

// Error logs at error level using default logger
func Error(format string, v ...interface{}) {
	Default().Error(format, v...)
}

// SetLevel sets the default logger level
func SetLevel(level Level) {
	Default().SetLevel(level)
}
