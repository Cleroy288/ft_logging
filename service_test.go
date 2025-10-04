package ft_logging

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

// TestResult tracks the outcome of a single test
type TestResult struct {
	Name   string
	Passed bool
	Output string
	Error  string
}

// Global test results tracker
var (
	testResults   []TestResult
	testResultsMu sync.Mutex
)

// recordTestResult records the result of a test
func recordTestResult(name string, passed bool, output string, err string) {
	testResultsMu.Lock()
	defer testResultsMu.Unlock()

	testResults = append(testResults, TestResult{
		Name:   name,
		Passed: passed,
		Output: output,
		Error:  err,
	})
}

// printTestSummary prints a colored summary of all test results
func printTestSummary() {
	var (
		totalTests  int
		passedTests int
		failedTests int
	)

	testResultsMu.Lock()
	defer testResultsMu.Unlock()

	totalTests = len(testResults)
	for _, result := range testResults {
		if result.Passed {
			passedTests++
		} else {
			failedTests++
		}
	}

	fmt.Println("\n\n" + string(bytes.Repeat([]byte("="), 60)))
	fmt.Println("                     TEST SUMMARY")
	fmt.Println(string(bytes.Repeat([]byte("="), 60)))

	// print individual test results
	for _, result := range testResults {
		if result.Passed {
			fmt.Printf("%s[SUCCESS]%s %s\n", colorGreen, colorReset, result.Name)
		} else {
			fmt.Printf("%s[FAIL]%s %s\n", colorRed, colorReset, result.Name)
			if result.Error != "" {
				fmt.Printf("  Error: %s\n", result.Error)
			}
			if result.Output != "" {
				fmt.Printf("  Output:\n%s\n", result.Output)
			}
		}
	}

	// print summary statistics
	fmt.Println(string(bytes.Repeat([]byte("="), 60)))
	fmt.Printf("Total:  %d tests\n", totalTests)
	fmt.Printf("%sPassed: %d tests%s\n", colorGreen, passedTests, colorReset)
	if failedTests > 0 {
		fmt.Printf("%sFailed: %d tests%s\n", colorRed, failedTests, colorReset)
	} else {
		fmt.Printf("Failed: %d tests\n", failedTests)
	}
	fmt.Println(string(bytes.Repeat([]byte("="), 60)))
}

// TestMain runs before all tests and prints summary after all tests complete
func TestMain(m *testing.M) {
	// run all tests
	exitCode := m.Run()

	// print summary
	printTestSummary()

	// exit with appropriate code
	os.Exit(exitCode)
}

func TestNewLogger(t *testing.T) {
	var (
		testName     = "TestNewLogger"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing NewLogger\n")
	output.WriteString("========================================\n")

	keys := []string{"request_id", "user_id"}
	logger := NewLogger(keys)
	if logger == nil {
		errorMessage = "NewLogger() returned nil"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Fatal(errorMessage)
		return
	}

	service, ok := logger.(*Service)
	if !ok {
		errorMessage = "NewLogger() did not return *Service"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Fatal(errorMessage)
		return
	}

	if len(service.contextKeys) != 2 {
		errorMessage = fmt.Sprintf("Expected 2 context keys, got %d", len(service.contextKeys))
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ NewLogger created successfully with context keys\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestNewLoggerNilKeys(t *testing.T) {
	var (
		testName     = "TestNewLoggerNilKeys"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing NewLoggerNilKeys\n")
	output.WriteString("========================================\n")

	logger := NewLogger(nil)
	service := logger.(*Service)

	if service.contextKeys != nil {
		errorMessage = fmt.Sprintf("Expected nil context keys, got %v", service.contextKeys)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ NewLogger handles nil context keys correctly\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestNewLoggerEmptyKeys(t *testing.T) {
	var (
		testName     = "TestNewLoggerEmptyKeys"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing NewLoggerEmptyKeys\n")
	output.WriteString("========================================\n")

	logger := NewLogger([]string{})
	service := logger.(*Service)

	if len(service.contextKeys) != 0 {
		errorMessage = fmt.Sprintf("Expected empty context keys, got %d", len(service.contextKeys))
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ NewLogger handles empty context keys correctly\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestInfo(t *testing.T) {
	var (
		testName     = "TestInfo"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing Info\n")
	output.WriteString("========================================\n")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	logger := NewLogger(nil)
	ctx := context.Background()

	logger.Info(ctx, "[TEST] test info message")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "[INFO]") {
		errorMessage = fmt.Sprintf("Expected [INFO] in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	if !strings.Contains(logOutput, "[TEST] test info message") {
		errorMessage = fmt.Sprintf("Expected '[TEST] test info message' in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Info() logs message correctly\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestSuccess(t *testing.T) {
	var (
		testName     = "TestSuccess"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing Success\n")
	output.WriteString("========================================\n")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	logger := NewLogger(nil)
	ctx := context.Background()

	logger.Success(ctx, "operation completed")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "[SUCCESS]") {
		errorMessage = fmt.Sprintf("Expected [SUCCESS] in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	if !strings.Contains(logOutput, "operation completed") {
		errorMessage = fmt.Sprintf("Expected 'operation completed' in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	if !strings.Contains(logOutput, colorGreen) {
		errorMessage = fmt.Sprintf("Expected green color code in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Success() logs message with green color\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestError(t *testing.T) {
	var (
		testName     = "TestError"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing Error\n")
	output.WriteString("========================================\n")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	logger := NewLogger(nil)
	ctx := context.Background()

	logger.Error(ctx, "something went wrong")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "[ERROR]") {
		errorMessage = fmt.Sprintf("Expected [ERROR] in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	if !strings.Contains(logOutput, "something went wrong") {
		errorMessage = fmt.Sprintf("Expected 'something went wrong' in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	if !strings.Contains(logOutput, colorRed) {
		errorMessage = fmt.Sprintf("Expected red color code in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Error() logs message with red color\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestColorCodes(t *testing.T) {
	var (
		testName     = "TestColorCodes"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing ColorCodes\n")
	output.WriteString("========================================\n")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	logger := NewLogger(nil)
	ctx := context.Background()

	// Test Info (white)
	buf.Reset()
	logger.Info(ctx, "white")
	if !strings.Contains(buf.String(), colorWhite) {
		errorMessage = "Info() should use white color"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Error(errorMessage)
		return
	}

	// Test Success (green)
	buf.Reset()
	logger.Success(ctx, "green")
	if !strings.Contains(buf.String(), colorGreen) {
		errorMessage = "Success() should use green color"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Error(errorMessage)
		return
	}

	// Test Error (red)
	buf.Reset()
	logger.Error(ctx, "red")
	if !strings.Contains(buf.String(), colorRed) {
		errorMessage = "Error() should use red color"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Error(errorMessage)
		return
	}

	output.WriteString("✓ All color codes working correctly\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestContextExtraction(t *testing.T) {
	var (
		testName     = "TestContextExtraction"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing ContextExtraction\n")
	output.WriteString("========================================\n")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	keys := []string{"request_id", "user_id", "trace_id"}
	logger := NewLogger(keys)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "abc123")
	ctx = context.WithValue(ctx, "user_id", "user-456")

	logger.Info(ctx, "test message")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "request_id=abc123") {
		errorMessage = fmt.Sprintf("Expected 'request_id=abc123' in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	if !strings.Contains(logOutput, "user_id=user-456") {
		errorMessage = fmt.Sprintf("Expected 'user_id=user-456' in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	// trace_id not in context, should not appear
	if strings.Contains(logOutput, "trace_id=") {
		errorMessage = fmt.Sprintf("Should not contain trace_id in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Context extraction working correctly\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestNoContextValues(t *testing.T) {
	var (
		testName     = "TestNoContextValues"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing NoContextValues\n")
	output.WriteString("========================================\n")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	keys := []string{"request_id", "user_id"}
	logger := NewLogger(keys)
	ctx := context.Background()

	logger.Info(ctx, "no context values")

	logOutput := buf.String()
	// Should not contain curly braces (no context info)
	if strings.Contains(logOutput, "{") {
		errorMessage = fmt.Sprintf("Expected no context info brackets, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	if !strings.Contains(logOutput, "no context values") {
		errorMessage = fmt.Sprintf("Expected message in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Handles empty context correctly\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestNilContext(t *testing.T) {
	var (
		testName     = "TestNilContext"
		output       bytes.Buffer
		errorMessage string
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing NilContext\n")
	output.WriteString("========================================\n")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	keys := []string{"request_id"}
	logger := NewLogger(keys)

	logger.Info(nil, "nil context")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "nil context") {
		errorMessage = fmt.Sprintf("Expected message in output, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}
	// Should not panic or contain context info
	if strings.Contains(logOutput, "{") {
		errorMessage = fmt.Sprintf("Expected no context info with nil context, got: %s", logOutput)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Handles nil context without panicking\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}
