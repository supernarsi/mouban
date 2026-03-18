# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
# Build
go build -o main main.go

# Run locally
go run main.go

# Production build (Linux amd64 binary + Docker)
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o main main.go
sudo docker build -t mythsman/mouban -f Dockerfile --platform=linux/amd64 .
```

## Test & Lint

```bash
# Run all tests
go test ./...

# Run single package tests
go test ./dao/
go test ./util/

# Run specific test
go test -run TestFunctionName ./path/to/package/

# Lint
golangci-lint run
```

## Architecture Overview

This is a Go web crawling backend service that extracts user ratings and annotations from Douban (books, movies, music, games). It serves as the backend for the hexo-douban project.

### Layer Structure

```
main.go (entry point + Gin routes)
    │
    ├── common/        # Initialization: config (Viper), database (GORM), logging (logrus)
    ├── controller/    # HTTP handlers: /admin/* and /guest/* endpoints
    ├── agent/         # Background workers: scheduled crawlers for users and items
    ├── crawl/         # Web scraping: HTTP client, HTML parsing, rate limiting
    ├── dao/           # Data access: GORM operations for all entities
    ├── model/         # GORM entities: Book, Movie, Game, Song, User, Comment, Rating, Schedule
    ├── consts/        # Constants: type codes, URL patterns
    └── util/          # Helpers: JSON parsing, stack traces
```

### Data Flow

1. **User Input** → `/guest/check_user?id={douban_id}` triggers user profile crawl
2. **User Crawl** → Fetches user overview + RSS feeds for wish/do/collect lists
3. **Comment Sync** → Scrapes user's comment pages, extracts item references
4. **Item Discovery** → New items queued to schedule table for detailed crawling
5. **Item Crawl** → Fetches full item details, discovers related users/items

### Key Patterns

- **Blank imports** for side effects: `_ "mouban/common"` triggers init functions
- **DAO layer** handles all database operations with GORM
- **Upsert pattern**: query existing, then update or create
- **Panic recovery** in middleware and goroutines with `defer recover()`
- **Environment overrides**: `KEY__SUBKEY` format maps to `key.subkey` in config

### Configuration

Uses Viper with YAML (`application.yml`) and environment variable overrides:

```bash
# Key env vars
GIN_MODE=release
server__cors=https://mydomain.com
datasource__host=localhost
datasource__username=root
datasource__password=secret
agent__enable=true
http__auth=dbcl2_cookie_value,http://user:pass@proxy:port;
```

### Database Schema

9 tables auto-migrated via GORM:
- `users` - User profiles with counts and timestamps
- `books/movies/games/songs` - Item metadata
- `comments` - User comments with ratings
- `ratings` - Aggregate rating data
- `schedules` - Crawler job queue with status
- `access` - Rate limiting tokens
- `storage` - Image/S3 cache references
