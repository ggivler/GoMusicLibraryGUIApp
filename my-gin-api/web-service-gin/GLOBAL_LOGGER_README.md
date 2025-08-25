# Global Logger Implementation

This document explains how the global logger has been implemented and how to use it throughout your Go application.

## Overview

The global logger allows you to write log messages from anywhere in your application without having to pass logger instances around or set up logging in each file.

## File Structure

```
├── main.go                    # Main application file using the global logger
├── logger/
│   └── logger.go             # Global logger package
├── example_usage.go          # Examples of how to use the logger
└── GLOBAL_LOGGER_README.md   # This documentation
```

## Logger Package (`logger/logger.go`)

The logger package provides:

- **Global Logger Instance**: A single logger instance accessible throughout the application
- **Convenient Functions**: Easy-to-use functions for different log levels
- **File Management**: Automatic file handling and cleanup

### Available Functions:

- `logger.Init(logFilePath)` - Initialize the global logger with a log file path
- `logger.Close()` - Close the log file (call this when shutting down)
- `logger.Info(v ...)` - Log info messages
- `logger.Infof(format, v ...)` - Log formatted info messages
- `logger.Error(v ...)` - Log error messages  
- `logger.Errorf(format, v ...)` - Log formatted error messages
- `logger.Debug(v ...)` - Log debug messages
- `logger.Debugf(format, v ...)` - Log formatted debug messages
- `logger.GetWriter()` - Get the underlying file writer (for frameworks like Gin)

## Usage Examples

### 1. Initialization (in main.go)

```go
func main() {
    // Initialize the global logger
    err := logger.Init("music-api.log")
    if err != nil {
        fmt.Printf("Failed to initialize logger: %v\\n", err)
        os.Exit(1)
    }
    defer logger.Close() // Ensure cleanup

    // Your application code here...
}
```

### 2. Using the Logger Anywhere

```go
package main

import "example/web-service-gin/logger"

func someFunction() {
    logger.Info("Function started")
    logger.Infof("Processing %d items", 42)
    
    if err := someOperation(); err != nil {
        logger.Errorf("Operation failed: %v", err)
        return
    }
    
    logger.Debug("Function completed successfully")
}
```

### 3. HTTP Handler Logging

```go
func getMusicFileInfo(c *gin.Context) {
    logger.Info("getMusicFileInfo endpoint called")
    // Handler logic...
    c.IndentedJSON(http.StatusOK, data)
}
```

### 4. Error Handling

```go
func processFile(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        logger.Errorf("Failed to read file %s: %v", filename, err)
        return err
    }
    
    logger.Infof("Successfully processed file: %s", filename)
    return nil
}
```

### 5. Integration with Gin Framework

```go
func main() {
    // ... logger initialization ...
    
    // Set Gin's output to use the same log file
    gin.DefaultWriter = logger.GetWriter()
    
    router := gin.Default()
    // Your routes...
}
```

## Benefits

1. **Centralized Logging**: All log messages go to the same file with consistent formatting
2. **Easy to Use**: Simple import and function calls from anywhere in your code
3. **No Parameter Passing**: No need to pass logger instances between functions
4. **Automatic Cleanup**: Built-in file handling and cleanup
5. **Framework Integration**: Easy integration with frameworks like Gin
6. **Different Log Levels**: Support for Info, Error, and Debug levels

## Log Format

The logger automatically includes:
- **Date and Time**: When the log entry was created
- **File and Line**: Where the log call was made from
- **Microseconds**: For precise timing
- **Log Level**: [ERROR] or [DEBUG] prefixes where applicable

Example log entry:
```
2025/08/24 21:02:24.461962 main.go:89: Application started.
2025/08/24 21:02:24.461962 main.go:90: Processing data for user: John Doe
2025/08/24 21:02:24.461962 logger.go:42: Starting server on localhost:8080
```

## Migration from Standard Log Package

If you were previously using Go's standard `log` package:

**Before:**
```go
log.Printf("Processing %d items", count)
log.Println("Operation completed")
```

**After:**
```go
logger.Infof("Processing %d items", count)
logger.Info("Operation completed")
```

## Best Practices

1. **Initialize Early**: Call `logger.Init()` as early as possible in your main function
2. **Use Appropriate Levels**: Use Info for general information, Error for errors, Debug for development
3. **Format Consistently**: Use formatted versions (`Infof`, `Errorf`) when you need variable substitution
4. **Clean Up**: Always use `defer logger.Close()` in your main function
5. **Import Consistently**: Use the full module path in imports: `"example/web-service-gin/logger"`

## Advanced Features

### Custom Log Levels
You can extend the logger package to add custom log levels like WARNING, FATAL, etc., by adding more functions similar to the existing ones.

### Multiple Log Files
If needed, you can modify the logger package to support writing to multiple files or different log levels to different files.

### Log Rotation
For production applications, consider adding log rotation functionality to prevent log files from growing too large.
