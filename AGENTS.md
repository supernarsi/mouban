# AGENTS.MD - Go Web Crawling Service Development Guidelines

## Project Overview
This is a Go-based web crawling backend service for Douban data (books, movies, music, games). It serves as the backend for the hexo-douban project to extract user ratings and annotations from Douban.

## Build Commands

### Standard Build 
```bash
go build -o main main.go
```

### Production Build with Docker
```bash
# Build Linux binary and Docker image
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o main main.go
sudo docker build -t mythsman/mouban -f Dockerfile --platform=linux/amd64 .
sudo docker push mythsman/mouban
```

### Local Run
```bash
# Using Go
go run main.go

# Or build and run
go build main.go && ./main
```

## Test Commands

### Run All Tests
```bash
go test ./...
```

### Run Specific Package Tests
```bash
go test ./util/
go test ./dao/
go test ./controller/
```

### Run Specific Test
```bash
go test -run TestFunctionName ./path/to/package/
```

### Run Single Test File
```bash
cd ./path/to/test/dir
go test -run TestFunctionName
```

### Test with Verbose Output
```bash
go test -v ./...
```

### Coverage Analysis
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Lint Commands
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
golangci-lint run --disable-all --enable=gofmt
```

## Code Style Guidelines

### Imports
- Group imports with blank lines separating std, third-party, and local packages
- Use blank import with comment for side effects: `_ "mouban/common"`
- Order: Standard library → Third-party → Local project packages

```go
import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"

    "mouban/consts"
    "mouban/controller"
    "mouban/util"
)
```

### Formatting
- Use `go fmt` for consistent formatting
- Line length should be reasonable (typically under 120 chars)
- Functions with long signatures may be formatted across multiple lines
- Struct fields aligned for readability

### Naming Conventions
- Use camelCase for function, variable, and method names
- Public items (exported) start with uppercase
- Private items start with lowercase
- Acronyms should be capitalized consistently: use `HTTP` not `Http`
- Constants: `UPPER_SNAKE_CASE` (if exported)

### Error Handling
- Handle panics using defer/recover in middleware and key functions 
- Log errors using logrus with appropriate log levels
- Use centralized error handling in gin middleware (see main.go handle function)
- Use structured logging with fields when possible
- For web responses, return structured JSON responses indicating success/failure

### Database Code Patterns
- Use DAO layer for database operations (data access objects)
- Use GORM ORM with struct tags like `gorm:"not null;uniqueIndex"`
- Follow upsert patterns: query existing record, then update or create
- Use meaningful variable names for DB operations (common.Db for connection)

### Struct Definition
- Define table name: `func (StructName) TableName() string`
- Use clear VO (View Object) structs for API responses
- Tag JSON fields for API outputs: `json:"field_name"`
- Use appropriate SQL tags for database columns: `gorm:"type:varchar(512);"`
- Follow DTO (Data Transfer Object) patterns with separate structs for database and API serialization

### Logging Pattern
- Use `logrus` for structured logging
- Use `logrus.Errorln()`, `logrus.Infoln()` for general messages
- Use `logrus.WithField().DebugInfo()` to add context to logs
- Standardize log message formatting across the codebase

### Web Routing
- Use Gin for routing with grouped routes
- Use consistent path structure: `/admin/...` for admin, `/guest/...` for guest
- Include comprehensive error handling in middleware
- Follow RESTful API principles with query parameters

### Testing Patterns
- Use table-driven tests with clear struct definitions for test cases
- Use subtests with `t.Run()` to namespace individual test cases
- Use proper test naming: `Test<FunctionName>`
- Include comprehensive test cases for different scenarios and edge cases
- Use `reflect.DeepEqual()` where appropriate for complex object comparison

### Configuration
- Use Viper for configuration management
- Keep config in application.yml (production) or use env vars
- Use consistent naming: `server.port`, `database.host`, etc.
- Use typed configuration values (GetInt, GetString, GetBool) 

### Concurrency and Error Recovery
- Use goroutines carefully with proper error handling
- Include panic recovery in goroutines with log output
- Use timer channels for periodic processes: `time.NewTicker()`
- Always handle potential panics in concurrent code with `defer recover()`
- Use `util.GetCurrentGoroutineStack()` for enhanced stack trace information in panic cases

### Enum and Type Definitions
- Define typed enums/structs for business constants (RatingStatus, Actions, Types)
- Use consistent patterns for enum validation with switch statements
- Follow parse functions with clear defaults (ParseType, ParseResult)

### Controller Patterns
- Return success/failure with structured JSON: `{"success": boolean, "msg": "error message"}`
- Use separate routes grouped by functionality (`/admin`, `/guest`)
- Implement proper request/response formatting

### Agent Architecture
- Implement background agent systems that run periodically
- Use fallback strategies and agent init() patterns
- Include orphan schedule cleanup and monitoring

## API Endpoints
- Main endpoints in main.go: check user, user items (book, movie, game, song)  
- Admin endpoints for load data, refresh item, refresh user
- All endpoints follow GET method pattern with query parameters
- Response format: {"result": data, "success": boolean}

## Environment Variables
Key env vars:
```
GIN_MODE=release
server__cors=https://mydomain.com
datasource__host=localhost
datasource__username=root
datasource__password=123456
```

## Development Notes
- Application handles web crawling with rate limiting to respect APIs
- Database uses MySQL with GORM for ORM layer  
- Uses Prometheus metrics for monitoring performance
- Implements rate limiting (token bucket pattern) and authentication handling for APIs
- Has structured data models for books, movies, music, games, users with appropriate relationships
- Implements agent system with scheduling for background processing
- Includes distributed crawling logic with proxy authentication support