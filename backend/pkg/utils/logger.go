package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger levels
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

// Log categories
const (
	CategoryAuth       = "AUTH"
	CategoryHandler    = "HANDLER"
	CategoryRoute      = "ROUTE"
	CategoryDatabase   = "DATABASE"
	CategoryIngestion  = "INGESTION"
	CategoryMigration  = "MIGRATION"
	CategoryGroq       = "Groq"
	CategoryQuery      = "QUERY"
	CategoryServer     = "SERVER"
	CategoryCache      = "CACHE"
	CategoryValidation = "VALIDATION"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string
	Level     string
	Category  string
	Message   string
	ClientIP  string
	Error     error
	Data      map[string]interface{}
}

// Logger is the main logging struct
type Logger struct {
	consoleLogger *log.Logger
	fileLogger    *log.Logger
	logFile       *os.File
	logDir        string
}

var globalLogger *Logger

// InitLogger initializes the global logger
func InitLogger(logsDir string) error {
	if logsDir == "" {
		logsDir = "logs"
	}

	// Create logs directory if it doesn't exist
	err := os.MkdirAll(logsDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create daily log file
	now := time.Now()
	logFileName := fmt.Sprintf("o2c-graph-%04d%02d%02d.log", now.Year(), now.Month(), now.Day())
	logFilePath := filepath.Join(logsDir, logFileName)

	// Open or create log file (append mode)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer for both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	globalLogger = &Logger{
		consoleLogger: log.New(multiWriter, "", 0),
		fileLogger:    log.New(logFile, "", 0),
		logFile:       logFile,
		logDir:        logsDir,
	}

	return nil
}

// GetLogger returns the global logger
func GetLogger() *Logger {
	if globalLogger == nil {
		// Fallback: initialize with default logs directory
		_ = InitLogger("logs")
	}
	return globalLogger
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// formatLog formats a log entry into a string
func (l *Logger) formatLog(entry LogEntry) string {
	// Build base log message
	var msg string
	if entry.ClientIP != "" {
		msg = fmt.Sprintf("[%s] [%s] [%s] [IP: %s] %s",
			entry.Timestamp,
			entry.Level,
			entry.Category,
			entry.ClientIP,
			entry.Message,
		)
	} else {
		msg = fmt.Sprintf("[%s] [%s] [%s] %s",
			entry.Timestamp,
			entry.Level,
			entry.Category,
			entry.Message,
		)
	}

	// Add error if present
	if entry.Error != nil {
		msg += fmt.Sprintf(" | Error: %v", entry.Error)
	}

	// Add data if present
	if len(entry.Data) > 0 {
		msg += fmt.Sprintf(" | Data: %+v", entry.Data)
	}

	return msg
}

// log is the internal logging method
func (l *Logger) log(level, category, message, clientIP string, err error, data map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000"),
		Level:     level,
		Category:  category,
		Message:   message,
		ClientIP:  clientIP,
		Error:     err,
		Data:      data,
	}

	formattedMsg := l.formatLog(entry)
	l.consoleLogger.Println(formattedMsg)
}

// Info logs an info message
func (l *Logger) Info(category, message string) {
	l.log(LevelInfo, category, message, "", nil, nil)
}

// InfoWithIP logs an info message with client IP
func (l *Logger) InfoWithIP(category, message, clientIP string) {
	l.log(LevelInfo, category, message, clientIP, nil, nil)
}

// InfoWithData logs an info message with additional data
func (l *Logger) InfoWithData(category, message string, data map[string]interface{}) {
	l.log(LevelInfo, category, message, "", nil, data)
}

// InfoWithDataIP logs an info message with data and client IP
func (l *Logger) InfoWithDataIP(category, message, clientIP string, data map[string]interface{}) {
	l.log(LevelInfo, category, message, clientIP, nil, data)
}

// Warn logs a warning message
func (l *Logger) Warn(category, message string) {
	l.log(LevelWarn, category, message, "", nil, nil)
}

// WarnWithIP logs a warning message with client IP
func (l *Logger) WarnWithIP(category, message, clientIP string) {
	l.log(LevelWarn, category, message, clientIP, nil, nil)
}

// WarnWithData logs a warning message with data
func (l *Logger) WarnWithData(category, message string, data map[string]interface{}) {
	l.log(LevelWarn, category, message, "", nil, data)
}

// WarnWithDataIP logs a warning message with data and client IP
func (l *Logger) WarnWithDataIP(category, message, clientIP string, data map[string]interface{}) {
	l.log(LevelWarn, category, message, clientIP, nil, data)
}

// Error logs an error message
func (l *Logger) Error(category, message string, err error) {
	l.log(LevelError, category, message, "", err, nil)
}

// ErrorWithIP logs an error message with client IP
func (l *Logger) ErrorWithIP(category, message, clientIP string, err error) {
	l.log(LevelError, category, message, clientIP, err, nil)
}

// ErrorWithData logs an error message with data
func (l *Logger) ErrorWithData(category, message string, err error, data map[string]interface{}) {
	l.log(LevelError, category, message, "", err, data)
}

// ErrorWithDataIP logs an error message with data and client IP
func (l *Logger) ErrorWithDataIP(category, message, clientIP string, err error, data map[string]interface{}) {
	l.log(LevelError, category, message, clientIP, err, data)
}

// Debug logs a debug message
func (l *Logger) Debug(category, message string) {
	l.log(LevelDebug, category, message, "", nil, nil)
}

// DebugWithIP logs a debug message with client IP
func (l *Logger) DebugWithIP(category, message, clientIP string) {
	l.log(LevelDebug, category, message, clientIP, nil, nil)
}

// DebugWithData logs a debug message with data
func (l *Logger) DebugWithData(category, message string, data map[string]interface{}) {
	l.log(LevelDebug, category, message, "", nil, data)
}

// DebugWithDataIP logs a debug message with data and client IP
func (l *Logger) DebugWithDataIP(category, message, clientIP string, data map[string]interface{}) {
	l.log(LevelDebug, category, message, clientIP, nil, data)
}

// ExtractClientIP extracts the client IP from request headers or remoteAddr
func ExtractClientIP(remoteAddr, xForwardedFor, xRealIP string) string {
	// Try X-Forwarded-For first (for proxies)
	if xForwardedFor != "" {
		return xForwardedFor
	}

	// Try X-Real-IP next
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to remote address
	if remoteAddr != "" {
		// Remove port from address (format: "IP:PORT")
		if idx := len(remoteAddr) - 1; idx > 0 && remoteAddr[idx-1:idx] == ":" {
			return remoteAddr[:idx-1]
		}
		return remoteAddr
	}

	return "unknown"
}

// Convenience functions using global logger

// Log logs with all parameters
func Log(level, category, message, clientIP string, err error, data map[string]interface{}) {
	GetLogger().log(level, category, message, clientIP, err, data)
}

// LogInfo logs an info message
func LogInfo(category, message string) {
	GetLogger().Info(category, message)
}

// LogInfoWithIP logs an info message with client IP
func LogInfoWithIP(category, message, clientIP string) {
	GetLogger().InfoWithIP(category, message, clientIP)
}

// LogInfoWithData logs an info message with data
func LogInfoWithData(category, message string, data map[string]interface{}) {
	GetLogger().InfoWithData(category, message, data)
}

// LogInfoWithDataIP logs an info message with data and IP
func LogInfoWithDataIP(category, message, clientIP string, data map[string]interface{}) {
	GetLogger().InfoWithDataIP(category, message, clientIP, data)
}

// LogWarn logs a warning message
func LogWarn(category, message string) {
	GetLogger().Warn(category, message)
}

// LogWarnWithIP logs a warning message with client IP
func LogWarnWithIP(category, message, clientIP string) {
	GetLogger().WarnWithIP(category, message, clientIP)
}

// LogWarnWithData logs a warning message with data
func LogWarnWithData(category, message string, data map[string]interface{}) {
	GetLogger().WarnWithData(category, message, data)
}

// LogWarnWithDataIP logs a warning message with data and IP
func LogWarnWithDataIP(category, message, clientIP string, data map[string]interface{}) {
	GetLogger().WarnWithDataIP(category, message, clientIP, data)
}

// LogError logs an error message
func LogError(category, message string, err error) {
	GetLogger().Error(category, message, err)
}

// LogErrorWithIP logs an error message with client IP
func LogErrorWithIP(category, message, clientIP string, err error) {
	GetLogger().ErrorWithIP(category, message, clientIP, err)
}

// LogErrorWithData logs an error message with data
func LogErrorWithData(category, message string, err error, data map[string]interface{}) {
	GetLogger().ErrorWithData(category, message, err, data)
}

// LogErrorWithDataIP logs an error message with data and IP
func LogErrorWithDataIP(category, message, clientIP string, err error, data map[string]interface{}) {
	GetLogger().ErrorWithDataIP(category, message, clientIP, err, data)
}

// LogDebug logs a debug message
func LogDebug(category, message string) {
	GetLogger().Debug(category, message)
}

// LogDebugWithIP logs a debug message with client IP
func LogDebugWithIP(category, message, clientIP string) {
	GetLogger().DebugWithIP(category, message, clientIP)
}

// LogDebugWithData logs a debug message with data
func LogDebugWithData(category, message string, data map[string]interface{}) {
	GetLogger().DebugWithData(category, message, data)
}

// LogDebugWithDataIP logs a debug message with data and IP
func LogDebugWithDataIP(category, message, clientIP string, data map[string]interface{}) {
	GetLogger().DebugWithDataIP(category, message, clientIP, data)
}
