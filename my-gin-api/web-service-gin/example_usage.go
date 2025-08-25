package main

import (
	"github.com/gin-gonic/gin"
	"example/web-service-gin/logger"
)

// ExampleFunction demonstrates how to use the global logger from any function
func ExampleFunction() {
	logger.Info("This is an info message from ExampleFunction")
	logger.Errorf("This is an error message with formatting: %s", "error details")
	logger.Debug("This is a debug message")
}

// ExampleMiddleware shows how you might use logging in a Gin middleware
func ExampleMiddleware() gin.HandlerFunc {
	return gin.LoggerWithWriter(logger.GetWriter())
}

// ExampleErrorHandler shows logging in error handling
func ExampleErrorHandler(err error) {
	if err != nil {
		logger.Errorf("An error occurred: %v", err)
		// Handle the error appropriately
	}
}
