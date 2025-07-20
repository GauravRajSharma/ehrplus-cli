# AGENTS.md - EHRPlus CLI Development Guide

## Build/Test Commands
- **Build:** `go build` or `just install` (builds and installs)
- **Run:** `go run main.go <command>`
- **Test:** `go test ./...` (standard Go testing)
- **Single test:** `go test -run TestName ./path/to/package`
- **Lint:** `go vet ./...` and `gofmt -s -w .`
- **Dependencies:** `go mod tidy` (clean up module dependencies)

## Code Style Guidelines
- **Package structure:** Use Cobra CLI pattern with `cmd/` directory for commands
- **Imports:** Group standard library, third-party, and local imports with blank lines
- **Naming:** Use camelCase for unexported, PascalCase for exported identifiers
- **Comments:** Use block comments for copyright headers, line comments for inline docs
- **Error handling:** Always handle errors explicitly, use `cobra.CheckErr()` for CLI errors
- **Types:** Define structs close to their usage, use typed constants for flags
- **TUI Models:** Name Bubble Tea models descriptively (e.g., `containerModel`, `demoModel`) to avoid conflicts

## Dependencies
- **CLI Framework:** Cobra for command structure, Viper for configuration
- **TUI Components:** Bubble Tea for terminal UI, Huh for forms, Lipgloss for styling, Bubbles for widgets
- **Database:** GORM ORM with SQLite driver for local data storage
- **SSH:** golang.org/x/crypto/ssh for remote deployment capabilities
- **Docker:** Docker client for container management

## Project Structure
- `main.go`: Entry point calling `cmd.Execute()`
- `cmd/`: Command definitions (root, system, version, demo commands)
- `cmd/demo.go`: Interactive TUI demonstration with forms, lists, and database examples
- `go.mod`: Uses Go 1.23.5 with TUI, database, and deployment dependencies
- `Justfile`: Simple build automation with `just install` command