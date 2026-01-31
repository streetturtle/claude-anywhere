# AGENTS.md - Developer Guide for Coding Agents

This guide provides essential information for AI coding agents working on the `claude-anywhere` project.

## Project Overview

`claude-anywhere` is a Go CLI tool that helps users browse and resume Claude Code sessions from anywhere. It reads session data from `~/.claude_sessions` and displays it in an interactive TUI powered by Bubbletea.

## Build, Test, and Lint Commands

### Build Commands
```bash
# Build the binary
make build
# Or: go build -o claude-anywhere ./cmd/claude-anywhere

# Run without building
make run
# Or: go run ./cmd/claude-anywhere

# Build and install to /usr/local/bin
make install

# Build for multiple platforms
make build-all
```

### Test Commands
```bash
# Run all tests
make test
# Or: go test ./...

# Run tests in a specific package
go test ./pkg/session

# Run a single test
go test -run TestName ./pkg/session

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Lint and Format Commands
```bash
# Format code (REQUIRED before committing)
gofmt -w .

# Simplify code and format
gofmt -s -w .

# Check formatting without modifying files
gofmt -d .

# List files that need formatting
gofmt -l .

# Run go vet (static analysis)
go vet ./...

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

### Clean Commands
```bash
make clean
# Or: rm -f claude-anywhere && go clean
```

## Code Style Guidelines

### Imports
- Use standard library imports first, then third-party imports, then internal imports
- Group imports with blank lines between categories
- Example:
```go
import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/charmbracelet/bubbletea"
    "github.com/spf13/cobra"

    "github.com/streetturtle/claude-anywhere/pkg/session"
)
```

### Formatting
- **CRITICAL**: Always run `gofmt -w .` before committing
- Use tabs for indentation (Go standard)
- Max line length: Keep reasonable (120 chars preferred, no hard limit)
- Use `gofmt -s` to simplify code where possible

### Types and Naming Conventions
- **Package names**: lowercase, single word (e.g., `session`, not `sessionManager`)
- **Exported types/functions**: PascalCase (e.g., `LoadSessions`, `Session`)
- **Unexported types/functions**: camelCase (e.g., `parseSessionFile`, `item`)
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE for groups (e.g., `StatusActive`, `ActivityThresholdMs`)
- **Interfaces**: Usually end with "er" suffix (e.g., `Reader`, `Writer`)
- **Acronyms**: Keep uppercase in names (e.g., `SessionID`, not `SessionId`; `CWD`, not `Cwd`)

### Error Handling
- Always check errors immediately after they occur
- Wrap errors with context using `fmt.Errorf("message: %w", err)` for error chain preservation
- Return errors rather than logging and continuing
- Use early returns to reduce nesting
- Example:
```go
func LoadSessions() ([]*Session, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("failed to get home directory: %w", err)
    }
    // ... continue
}
```

### Comments
- Use `//` for single-line comments
- Exported functions, types, and constants MUST have doc comments
- Doc comments should start with the name of the thing being documented
- Example:
```go
// Session represents a Claude Code session with all its metadata
type Session struct {
    // ...
}

// LoadSessions reads all claude-status*.json files from ~/.claude_sessions and parses them
func LoadSessions() ([]*Session, error) {
    // ...
}
```

### Project Structure
```
claude-anywhere/
├── cmd/claude-anywhere/  # Main application entry point
│   └── main.go           # CLI commands using cobra
├── pkg/session/          # Core session logic
│   ├── session.go        # Session data structures and loading
│   └── tui.go            # Terminal UI using bubbletea
├── assets/               # Shell scripts and assets
├── Makefile              # Build automation
└── go.mod                # Go module definition
```

### Dependencies
- **CLI Framework**: `github.com/spf13/cobra` - command-line interface
- **TUI Framework**: `github.com/charmbracelet/bubbletea` - terminal UI
- **TUI Components**: `github.com/charmbracelet/bubbles` - reusable TUI components
- **Styling**: `github.com/charmbracelet/lipgloss` - terminal styling

### JSON Handling
- Use struct tags for JSON marshaling/unmarshaling
- Use `json:"-"` to exclude fields from JSON
- Handle empty/missing fields gracefully
- Example:
```go
type Session struct {
    SessionID      string    `json:"session_id"`
    StatusFilePath string    `json:"-"` // Not included in JSON
}
```

### Constants and Magic Numbers
- Define constants for thresholds and magic numbers
- Group related constants together
- Use comments to explain units (e.g., milliseconds)
- Example:
```go
const (
    ActivityThresholdMs = 3000      // 3 seconds
    ClosedThresholdMs   = 3600000   // 1 hour
    MaxAgeMs            = 604800000 // 7 days
)
```

### Testing (when adding tests)
- Test files should be named `*_test.go`
- Test functions should start with `Test`
- Use table-driven tests for multiple test cases
- Use `t.Helper()` for test helper functions
- Example:
```go
func TestLoadSessions(t *testing.T) {
    sessions, err := LoadSessions()
    if err != nil {
        t.Fatalf("LoadSessions failed: %v", err)
    }
    // ...
}
```

## Common Tasks

### Adding a new session field
1. Add field to `Session` struct in `pkg/session/session.go`
2. Add corresponding field to `StatusLineData` struct
3. Update `parseSessionFile()` to map the data
4. Update TUI rendering in `pkg/session/tui.go` if needed
5. Run `gofmt -w .` and test

### Modifying the TUI
- Edit `pkg/session/tui.go`
- Use lipgloss styles defined at package level
- Follow the Elm architecture (Model, Update, View)
- Test with `make run`

### Adding a new CLI command
- Edit `cmd/claude-anywhere/main.go`
- Create a cobra.Command
- Add it to rootCmd with `rootCmd.AddCommand()`
- Follow existing pattern (see `listCmd`)

## Development Workflow

1. Make changes to code
2. Format code: `gofmt -w .`
3. Run vet: `go vet ./...`
4. Test locally: `make run`
5. Build: `make build`
6. Test the binary: `./claude-anywhere`
7. Commit changes with descriptive message

## Important Notes

- **No test files yet**: The project doesn't have tests. When adding tests, follow Go testing conventions.
- **No linter config**: Use standard Go tools (`gofmt`, `go vet`). Consider adding `golangci-lint` in the future.
- **Home directory**: Sessions are always read from `~/.claude_sessions`
- **Session resumption**: Uses `claude -r <session-id>` command
- **TUI navigation**: Uses bubbletea's built-in list component with keyboard controls
