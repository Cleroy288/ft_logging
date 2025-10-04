// Package ft_logging provides simple, colorized logging with context support.
// It offers three logging levels: Info (white), Success (green), and Error (red).
package ft_logging

import (
	"context"
	"fmt"
	"log"
)

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorWhite  = "\033[37m"
)

// Logger defines the logging interface with three methods.
type Logger interface {
	// Info logs an informational message in white.
	Info(ctx context.Context, message string)

	// Success logs a success message in green.
	Success(ctx context.Context, message string)

	// Error logs an error message in red.
	Error(ctx context.Context, message string)
}

// Service implements the Logger interface.
type Service struct {
	contextKeys []string
}

// NewLogger creates a new Logger instance with optional context keys to extract.
// contextKeys is a slice of context keys to extract and log (pass nil or empty slice if not needed).
// Returns a Logger interface implementation.
//
// Example:
//   logger := ft_logging.NewLogger([]string{"request_id", "user_id", "trace_id"})
func NewLogger(contextKeys []string) Logger {
	var (
		service *Service
		keys    string
		i       int
		key     string
	)

	service = &Service{
		contextKeys: contextKeys,
	}

	// print initialization details
	if len(contextKeys) == 0 {
		log.Printf("[ft_logging] Initialized with no context extraction")
	} else {
		keys = ""
		for i, key = range contextKeys {
			if i > 0 {
				keys += ", "
			}
			keys += key
		}
		log.Printf("[ft_logging] Initialized with context keys: [%s]", keys)
	}

	return service
}

// Info logs an informational message in white.
// ctx is the context for extracting context values.
// message is the log message to display.
func (s *Service) Info(ctx context.Context, message string) {
	s.LogWithColor(ctx, colorWhite, "INFO", message)
}

// Success logs a success message in green.
// ctx is the context for extracting context values.
// message is the log message to display.
func (s *Service) Success(ctx context.Context, message string) {
	s.LogWithColor(ctx, colorGreen, "SUCCESS", message)
}

// Error logs an error message in red.
// ctx is the context for extracting context values.
// message is the log message to display.
func (s *Service) Error(ctx context.Context, message string) {
	s.LogWithColor(ctx, colorRed, "ERROR", message)
}

// LogWithColor formats and logs messages with color and context information.
// ctx is the context for extracting context values.
// color is the ANSI color code for the log level.
// level is the log level name (INFO, SUCCESS, ERROR).
// message is the log message to display.
func (s *Service) LogWithColor(ctx context.Context, color, level, message string) {
	var (
		contextInfo  string
		contextPart  string
		formattedMsg string
	)

	// extract context information
	contextInfo = s.extractContextInfo(ctx)
	contextPart = ""
	if contextInfo != "" {
		contextPart = fmt.Sprintf(" {%s}", contextInfo)
	}

	formattedMsg = fmt.Sprintf("%s[%s]%s %s%s", color, level, colorReset, message, contextPart)
	log.Print(formattedMsg)
}

// extractContextInfo extracts context values using the configured keys.
// ctx is the context to extract values from.
// Returns a formatted string with all extracted context values.
func (s *Service) extractContextInfo(ctx context.Context) string {
	var (
		parts  []string
		key    string
		value  any
		result string
		i      int
		part   string
	)

	if ctx == nil || len(s.contextKeys) == 0 {
		return ""
	}

	parts = []string{}

	// loop through configured context keys
	for _, key = range s.contextKeys {
		value = ctx.Value(key)
		if value != nil {
			parts = append(parts, fmt.Sprintf("%s=%v", key, value))
		}
	}

	if len(parts) == 0 {
		return ""
	}

	result = ""
	for i, part = range parts {
		if i > 0 {
			result += ", "
		}
		result += part
	}
	return result
}
