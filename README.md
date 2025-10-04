# ft_logging

**Simple, colorized logging for Go.** Three functions. Three colors. Context extraction. Clean code.

## Features

- ✅ **Simple interface** - Only 3 methods: Info, Success, Error
- ✅ **Colorized output** - White (Info), Green (Success), Red (Error)
- ✅ **Context-aware** - Automatically extracts and logs context values
- ✅ **Configurable context keys** - Define which context keys to extract
- ✅ **Flexible formatting** - Add your own prefixes in messages
- ✅ **Standard library** - Uses Go's built-in log package
- ✅ **Zero dependencies** - No external packages required
- ✅ **Production-ready** - 100% test coverage

## Installation

```bash
go get github.com/Cleroy288/ft_logging
```

## Quick Start

```go
package main

import (
    "context"
    ft_logging "github.com/Cleroy288/ft_logging"
)

func main() {
    // Define context keys to extract
    keys := []string{"request_id", "user_id", "trace_id"}
    logger := ft_logging.NewLogger(keys)
    // Prints: [ft_logging] Initialized with context keys: [request_id, user_id, trace_id]

    // Add values to context
    ctx := context.Background()
    ctx = context.WithValue(ctx, "request_id", "abc123")
    ctx = context.WithValue(ctx, "user_id", "user-456")

    // Log messages (add your own prefixes)
    logger.Info(ctx, "[API] Processing request")
    logger.Success(ctx, "[DB] User authenticated")
    logger.Error(ctx, "[CACHE] Connection failed")
}
```

**Output:**
```
[ft_logging] Initialized with context keys: [request_id, user_id, trace_id]
[INFO] [API] Processing request {request_id=abc123, user_id=user-456}
[SUCCESS] [DB] User authenticated {request_id=abc123, user_id=user-456}
[ERROR] [CACHE] Connection failed {request_id=abc123, user_id=user-456}
```

## API Reference

### `NewLogger(contextKeys []string) Logger`

Creates a new Logger instance with optional context keys to extract. Prints initialization status on creation.

**Parameters:**
- `contextKeys` - Slice of context keys to extract from context. Pass `nil` or empty slice for no context extraction.

**Returns:**
- `Logger` - Logger interface implementation

**Initialization Output:**
```go
// With context extraction
logger := ft_logging.NewLogger([]string{"request_id", "user_id", "trace_id"})
// Prints: [ft_logging] Initialized with context keys: [request_id, user_id, trace_id]

// Without context extraction
logger := ft_logging.NewLogger(nil)
// Prints: [ft_logging] Initialized with no context extraction
```

This helps verify the logger was initialized correctly with the expected configuration.

### `Info(ctx context.Context, message string)`

Logs an informational message in **white**.

**Example:**
```go
logger.Info(ctx, "Processing request")
```

### `Success(ctx context.Context, message string)`

Logs a success message in **green**.

**Example:**
```go
logger.Success(ctx, "User created successfully")
```

### `Error(ctx context.Context, message string)`

Logs an error message in **red**.

**Example:**
```go
logger.Error(ctx, "Database connection failed")
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "context"
    ft_logging "github.com/Cleroy288/ft_logging"
)

func main() {
    // Simple logger without context extraction
    logger := ft_logging.NewLogger(nil)
    ctx := context.Background()

    logger.Info(ctx, "[APP] Starting server on port 8080")
    logger.Success(ctx, "[APP] Server started successfully")
    logger.Error(ctx, "[APP] Could not connect to database")
}
```

### With Context Extraction

```go
package main

import (
    "context"
    ft_logging "github.com/Cleroy288/ft_logging"
)

func main() {
    // Logger with context key extraction
    keys := []string{"request_id", "user_id", "session_id"}
    logger := ft_logging.NewLogger(keys)

    // Add values to context
    ctx := context.Background()
    ctx = context.WithValue(ctx, "request_id", "req-123")
    ctx = context.WithValue(ctx, "user_id", "user-789")

    logger.Info(ctx, "[AUTH] User login attempt")
    // Output: [INFO] [AUTH] User login attempt {request_id=req-123, user_id=user-789}
}
```

### In HTTP Handlers

```go
func handleRequest(logger ft_logging.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Add request metadata to context
        requestID := generateRequestID()
        ctx = context.WithValue(ctx, "request_id", requestID)
        ctx = context.WithValue(ctx, "method", r.Method)
        ctx = context.WithValue(ctx, "path", r.URL.Path)

        logger.Info(ctx, "[HTTP] Received request")

        if err := processRequest(r); err != nil {
            logger.Error(ctx, "[HTTP] Request failed: " + err.Error())
            http.Error(w, "Internal Server Error", 500)
            return
        }

        logger.Success(ctx, "[HTTP] Request completed")
        w.WriteHeader(http.StatusOK)
    }
}

// Create logger with HTTP context keys
logger := ft_logging.NewLogger([]string{"request_id", "method", "path", "user_id"})
```

### Dependency Injection

```go
type Service struct {
    logger ft_logging.Logger
}

func NewService(logger ft_logging.Logger) *Service {
    return &Service{logger: logger}
}

func (s *Service) DoWork(ctx context.Context) error {
    s.logger.Info(ctx, "Starting work")

    if err := s.performTask(); err != nil {
        s.logger.Error(ctx, "Task failed: " + err.Error())
        return err
    }

    s.logger.Success(ctx, "Work completed")
    return nil
}
```

### Different Components with Custom Prefixes

```go
// Different loggers for different components
apiLogger := ft_logging.NewLogger([]string{"request_id", "user_id"})
dbLogger := ft_logging.NewLogger([]string{"query_id", "table"})
cacheLogger := ft_logging.NewLogger([]string{"key", "operation"})

// Use them with custom prefixes in messages
apiLogger.Info(ctx, "[API] Processing request")
dbLogger.Info(ctx, "[DATABASE] Executing query")
cacheLogger.Success(ctx, "[CACHE] Key retrieved")
```

## Color Output

The package uses ANSI color codes for terminal output:

- **Info** - White (`\033[37m`) - For general information
- **Success** - Green (`\033[32m`) - For successful operations
- **Error** - Red (`\033[31m`) - For errors and failures

Colors are automatically reset after each message.

## Interface

```go
type Logger interface {
    Info(ctx context.Context, message string)
    Success(ctx context.Context, message string)
    Error(ctx context.Context, message string)
}
```

This interface makes it easy to:
- Mock logging in tests
- Swap implementations
- Use dependency injection

## Testing

Run all tests:

```bash
go test -v
```

Run with coverage:

```bash
go test -v -cover
```

**Test coverage:** 100% ✅

## Project Structure

```
ft_logging/
├── service.go       # Main implementation
├── service_test.go  # Comprehensive tests
├── go.mod           # Module definition
└── README.md        # This file
```

## Best Practices

### 1. Define context keys at initialization

```go
// Define all context keys you want to track
keys := []string{
    "request_id",
    "user_id",
    "trace_id",
    "session_id",
    "correlation_id",
}
logger := ft_logging.NewLogger(keys)
```

### 2. Add custom prefixes in messages

```go
// Use prefixes to identify components
logger.Info(ctx, "[API] Processing request")
logger.Success(ctx, "[DATABASE] Query executed")
logger.Error(ctx, "[CACHE] Connection failed")
```

### 3. Populate context in middleware

```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        ctx = context.WithValue(ctx, "request_id", generateID())
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 4. Use consistent context keys

```go
// Good: Consistent key naming
ctx = context.WithValue(ctx, "request_id", "abc123")
ctx = context.WithValue(ctx, "user_id", "user-456")

// Bad: Inconsistent naming
ctx = context.WithValue(ctx, "RequestID", "abc123")  // Won't match
ctx = context.WithValue(ctx, "uid", "user-456")      // Won't match
```

### 5. Clear, actionable messages

```go
// Good: Specific, informative
logger.Error(ctx, "[AUTH] Failed to verify JWT: token expired")

// Bad: Vague
logger.Error(ctx, "error")
```

## Comparison with Other Packages

| Feature | ft_logging | logrus | zap | zerolog |
|---------|-----------|--------|-----|---------|
| Simplicity | **Excellent** | Medium | Complex | Medium |
| Colors | ✅ | ✅ | ❌ | ❌ |
| Dependencies | **0** | Many | Some | 0 |
| Context support | ✅ | ✅ | ✅ | ✅ |
| Learning curve | **Easy** | Medium | Hard | Medium |
| Performance | Good | Good | **Excellent** | **Excellent** |

## Why ft_logging?

- **Simplicity** - 3 functions, minimal configuration
- **Readability** - Colored output makes logs easy to scan
- **Flexible** - Define your own context keys and message formats
- **Context extraction** - Automatic extraction of request IDs, user IDs, etc.
- **Standard** - Uses Go's built-in log package
- **Lightweight** - No external dependencies
- **Distributed tracing ready** - Track requests across services

Perfect for:
- Microservices with distributed tracing
- HTTP APIs needing request tracking
- Applications using context propagation
- Small to medium applications
- Prototypes and MVPs
- Services where simplicity matters

## Author

**Charles Leroy** - [@Cleroy288](https://github.com/Cleroy288)

---

**Keep it simple. Keep it clean.**
