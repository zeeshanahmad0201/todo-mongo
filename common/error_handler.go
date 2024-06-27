package common

import (
	"context"
	"log"
	"os"
	"runtime/debug"
	"time"
)

// ErrorHandlerConfig holds the configuration for the error handler.
type ErrorHandlerConfig struct {
	Message         string
	Exit            bool
	PrintStackTrace bool
}

// DefaultErrorHandlerConfig provides default values for error handler configuration.
var DefaultErrorHandlerConfig = ErrorHandlerConfig{
	Message:         "",
	Exit:            false,
	PrintStackTrace: true,
}

// HandleError logs the error with a custom message and decides whether to exit the program.
// The function can optionally print the stack trace for better debugging.
func HandleError(err error, config ...ErrorHandlerConfig) {
	if err != nil {
		// Use the provided configuration or the default one.
		cfg := DefaultErrorHandlerConfig
		if len(config) > 0 {
			cfg = config[0]
		}

		// Log the error with the custom message or a default error message.
		if cfg.Message != "" {
			log.Printf("%s: %v", cfg.Message, err)
		} else {
			log.Printf("Error: %v", err)
		}

		// Optionally print the stack trace for debugging purposes.
		if cfg.PrintStackTrace {
			log.Println(string(debug.Stack()))
		}

		// Exit the program if the configuration requires it.
		if cfg.Exit {
			os.Exit(1)
		}
	}
}

// creates a context with a timeout
func CreateContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
