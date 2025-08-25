package logger

import (
	"log"
	"os"
)

var (
	// Logger is the global logger instance
	Logger *log.Logger
	// logFile holds the file handle for proper cleanup
	logFile *os.File
)

// Init initializes the global logger with the specified log file
func Init(logFilePath string) error {
	// Open or create the log file with append, create, and write-only flags
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	
	logFile = file
	
	// Create a new logger instance
	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	
	return nil
}

// Close closes the log file (should be called when the application shuts down)
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// Info logs an info message
func Info(v ...interface{}) {
	if Logger != nil {
		Logger.Println(v...)
	}
}

// Infof logs a formatted info message
func Infof(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf(format, v...)
	}
}

// Error logs an error message
func Error(v ...interface{}) {
	if Logger != nil {
		Logger.Println("[ERROR]", v)
	}
}

// Errorf logs a formatted error message
func Errorf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[ERROR] "+format, v...)
	}
}

// Debug logs a debug message
func Debug(v ...interface{}) {
	if Logger != nil {
		Logger.Println("[DEBUG]", v)
	}
}

// Debugf logs a formatted debug message
func Debugf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[DEBUG] "+format, v...)
	}
}

// GetWriter returns the underlying writer for use with Gin or other frameworks
func GetWriter() *os.File {
	return logFile
}
