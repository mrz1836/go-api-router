# CLAUDE.md

Quick reference for Claude Code to work effectively with the go-api-router project.

## 🎯 Project Overview

**go-api-router** is a lightweight HTTP middleware library built on Julien Schmidt's httprouter, providing CORS handling, structured logging, JWT authentication, and standardized API responses.

## 📁 Key Components

| File                 | Purpose                                                          |
|----------------------|------------------------------------------------------------------|
| `api_router.go`      | Core router with middleware, CORS, logging, and request handling |
| `authentication.go`  | JWT token creation, validation, and management                   |
| `middleware.go`      | Middleware stack interface and standard HTTP adapter             |
| `response.go`        | JSON response helpers and structured API responses               |
| `error.go`           | Standardized API error handling and logging                      |
| `utilities.go`       | Helper functions (IP detection, parameter filtering, caching)    |
| `response_writer.go` | Custom response writer with status tracking                      |

## 🏗️ Architecture

```
Request → Router.Request() → Middleware Stack → Handler → APIResponseWriter → Response
              ↓                                              ↓
          Logging/CORS                                Status Tracking
```

- **Router**: Main entry point with configurable CORS, logging, and authentication
- **Middleware**: Chainable request processors using `httprouter.Handle` signature
- **Authentication**: JWT-based with configurable expiration and validation
- **Logging**: Structured logs with request IDs, timing, and parameter filtering

## 🛠️ Development Commands

```bash
# Build & Test
magex test              # Run standard test suite
magex test:race         # Run with race detector
magex test:cover        # Generate coverage report
magex bench             # Run benchmarks
magex lint              # Code quality checks

# Build
magex build             # Build for current platform
magex install           # Install to $GOPATH/bin

# Dependencies
magex deps:update       # Update all dependencies
magex tidy              # Run go mod tidy
```

## 🧪 Testing Standards

- **Coverage**: Maintain high test coverage for all public APIs
- **Race Detection**: Always run `magex test:race` before commits
- **Fuzz Tests**: Located in `*_fuzz_test.go` files
- **Test Structure**: Use `testify/assert` and `testify/require`
- **Parallel Tests**: Use `t.Parallel()` in unit tests

## 🔑 Key Patterns

### Router Configuration
```go
router := apirouter.New()
router.CrossOriginEnabled = true
router.FilterFields = []string{"password", "token"}
router.HTTPRouter.GET("/api/v1/user", router.Request(handler))
```

### Middleware Usage
```go
stack := apirouter.NewStack()
stack.Use(authMiddleware)
stack.Use(loggingMiddleware)
handler := stack.Wrap(finalHandler)
```

### Error Handling
```go
return apirouter.ErrorFromResponse(w, "internal message", "public message",
    errorCode, http.StatusBadRequest, data)
```

### Authentication
```go
authenticated, req, err := apirouter.Check(w, r, sessionSecret, issuer, sessionAge)
claims := apirouter.GetClaims(req)
```

## 📋 Common Tasks

**Adding New Middleware:**
1. Implement `func(httprouter.Handle) httprouter.Handle` signature
2. Add to middleware stack with `stack.Use(middleware)`
3. Test with existing router patterns

**Adding Authentication to Route:**
1. Use `router.Request()` for logged routes
2. Call `apirouter.Check()` in handler
3. Access user data via `apirouter.GetClaims(req)`

**Error Responses:**
1. Use `apirouter.RespondWith()` for JSON responses
2. Use `apirouter.ErrorFromResponse()` for structured errors
3. Sensitive fields auto-filtered from logs via `FilterFields`

**Testing New Features:**
1. Write unit tests with `testify`
2. Add fuzz tests for public APIs handling external input
3. Include race condition tests for concurrent access
4. Run full test suite: `magex test:full`

## 🔧 Code Quality

- **Linting**: Uses golangci-lint via `magex lint`
- **Formatting**: Standard Go formatting with `go fmt`
- **Documentation**: Godoc comments for all public APIs
- **Dependencies**: Security scanning via `magex deps:audit`

## 🚀 Release Process

1. Update version: `magex version:bump bump=patch push`
2. Validate: `magex release:test`
3. Create release: `magex release`

## 💡 Tips for Claude Code

- Always run tests before making changes: `magex test`
- Check existing patterns in test files before implementing new features
- Use `router.Request()` wrapper for routes needing logging/CORS
- JWT tokens automatically refreshed on valid requests
- IP addresses safely extracted behind load balancers via `GetClientIPAddress()`
- Parameters auto-parsed and cached via `GetParams()`
- Sensitive fields filtered from logs using `FilterFields` configuration
