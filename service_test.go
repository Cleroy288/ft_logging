package ft_logging

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	keys := []string{"request_id", "user_id"}
	logger := NewLogger(keys)
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	service, ok := logger.(*Service)
	if !ok {
		t.Fatal("NewLogger() did not return *Service")
	}

	if len(service.contextKeys) != 2 {
		t.Errorf("Expected 2 context keys, got %d", len(service.contextKeys))
	}
}

func TestNewLoggerNilKeys(t *testing.T) {
	logger := NewLogger(nil)
	service := logger.(*Service)

	if service.contextKeys != nil {
		t.Errorf("Expected nil context keys, got %v", service.contextKeys)
	}
}

func TestNewLoggerEmptyKeys(t *testing.T) {
	logger := NewLogger([]string{})
	service := logger.(*Service)

	if len(service.contextKeys) != 0 {
		t.Errorf("Expected empty context keys, got %d", len(service.contextKeys))
	}
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	logger := NewLogger(nil)
	ctx := context.Background()

	logger.Info(ctx, "[TEST] test info message")

	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Expected [INFO] in output, got: %s", output)
	}
	if !strings.Contains(output, "[TEST] test info message") {
		t.Errorf("Expected '[TEST] test info message' in output, got: %s", output)
	}
}

func TestSuccess(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	logger := NewLogger(nil)
	ctx := context.Background()

	logger.Success(ctx, "operation completed")

	output := buf.String()
	if !strings.Contains(output, "[SUCCESS]") {
		t.Errorf("Expected [SUCCESS] in output, got: %s", output)
	}
	if !strings.Contains(output, "operation completed") {
		t.Errorf("Expected 'operation completed' in output, got: %s", output)
	}
	if !strings.Contains(output, colorGreen) {
		t.Errorf("Expected green color code in output, got: %s", output)
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	logger := NewLogger(nil)
	ctx := context.Background()

	logger.Error(ctx, "something went wrong")

	output := buf.String()
	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("Expected [ERROR] in output, got: %s", output)
	}
	if !strings.Contains(output, "something went wrong") {
		t.Errorf("Expected 'something went wrong' in output, got: %s", output)
	}
	if !strings.Contains(output, colorRed) {
		t.Errorf("Expected red color code in output, got: %s", output)
	}
}

func TestColorCodes(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	logger := NewLogger(nil)
	ctx := context.Background()

	// Test Info (white)
	buf.Reset()
	logger.Info(ctx, "white")
	if !strings.Contains(buf.String(), colorWhite) {
		t.Error("Info() should use white color")
	}

	// Test Success (green)
	buf.Reset()
	logger.Success(ctx, "green")
	if !strings.Contains(buf.String(), colorGreen) {
		t.Error("Success() should use green color")
	}

	// Test Error (red)
	buf.Reset()
	logger.Error(ctx, "red")
	if !strings.Contains(buf.String(), colorRed) {
		t.Error("Error() should use red color")
	}
}

func TestContextExtraction(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	keys := []string{"request_id", "user_id", "trace_id"}
	logger := NewLogger(keys)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "abc123")
	ctx = context.WithValue(ctx, "user_id", "user-456")

	logger.Info(ctx, "test message")

	output := buf.String()
	if !strings.Contains(output, "request_id=abc123") {
		t.Errorf("Expected 'request_id=abc123' in output, got: %s", output)
	}
	if !strings.Contains(output, "user_id=user-456") {
		t.Errorf("Expected 'user_id=user-456' in output, got: %s", output)
	}
	// trace_id not in context, should not appear
	if strings.Contains(output, "trace_id=") {
		t.Errorf("Should not contain trace_id in output, got: %s", output)
	}
}

func TestNoContextValues(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	keys := []string{"request_id", "user_id"}
	logger := NewLogger(keys)
	ctx := context.Background()

	logger.Info(ctx, "no context values")

	output := buf.String()
	// Should not contain curly braces (no context info)
	if strings.Contains(output, "{") {
		t.Errorf("Expected no context info brackets, got: %s", output)
	}
	if !strings.Contains(output, "no context values") {
		t.Errorf("Expected message in output, got: %s", output)
	}
}

func TestNilContext(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	keys := []string{"request_id"}
	logger := NewLogger(keys)

	logger.Info(nil, "nil context")

	output := buf.String()
	if !strings.Contains(output, "nil context") {
		t.Errorf("Expected message in output, got: %s", output)
	}
	// Should not panic or contain context info
	if strings.Contains(output, "{") {
		t.Errorf("Expected no context info with nil context, got: %s", output)
	}
}
