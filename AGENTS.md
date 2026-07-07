## Development Guidelines

Important: Write tests for any changes you make and any test coverage gaps also.

## Go Coding Standards & Idioms
When writing or refactoring Go code in this repository, strictly adhere to the following rules:

1. **Error Handling:** - Never use `panic()` unless explicitly initializing a package-level variable that must not fail (e.g., compiling a regex). 
   - Always check errors explicitly (`if err != nil`). Do not ignore errors with `_`.
   - Wrap errors with context using `fmt.Errorf("doing something: %w", err)` to maintain the stack trace.

2. **Concurrency:**
   - Always pass `context.Context` as the first argument to functions that do I/O or heavy processing.
   - Guard goroutines with a `sync.WaitGroup` or an `errgroup.Group`. Do not leave orphaned goroutines.
   - Keep channel operations clean and ensure channels are closed by the sender, not the receiver.

3. **Style & Formatting:**
   - Follow [Effective Go](https://go.dev/doc/effective_go).
   - Use short variable names for limited scopes (e.g., `r` for `http.Request`, `db` for `*sql.DB`).
   - Keep interfaces small, typically 1-3 methods, defined where they are consumed, not where they are implemented.
   - Avoid package-level state/globals. Use dependency injection by passing structs with interfaces.

4. **Testing:**
   - Write table-driven tests for all business logic.
   - Use the standard `testing` package. Avoid third-party assertion libraries (like `testify`) unless they are already present in the specific package you are editing.
   - Name test files `[filename]_test.go`.

### Code Quality
- Write clean, idiomatic Go code following Go conventions and best practices
- Use structured logging with logrus for consistent log formatting
- Ensure proper error handling and meaningful error messages
- Follow the existing code style and patterns in the repository

### Testing and Quality Assurance
- **CRITICAL**: Always run ALL of the following commands before making a commit or opening a PR:
  1. `go fmt ./...` - Format all Go files
  2. `golangci-lint run` - Run all configured linters and formatters
  3. `Run all unit tests
- Ensure ALL tests pass AND ALL linting checks pass before committing
- The project uses golangci-lint with strict formatting rules - code must pass ALL checks

## Formatting and Linting Requirements

This project uses golangci-lint with strict formatting rules configured in `.golangci.yml`. Common formatting requirements include:

- Proper spacing around operators (`if condition {` not `if(condition){`)
- Correct struct field alignment and spacing
- Proper import ordering (standard library, third-party, local packages)
- No trailing whitespace
- Consistent spacing around assignment operators (`key: value` not `key:value`)
- Space after commas in function parameters and struct literals

**Always run `go fmt ./...`, `golangci-lint run`, AND `make test` after making ANY code changes to ensure both functionality and formatting are correct before committing.**

## Key File References

- Main entry: `cmd/podsync/main.go`
- Config loading: `cmd/podsync/config.go`
- Feed update: `services/update/updater.go` (episode lifecycle, cleanup)
- Episode filtering: `services/update/matcher.go`
- Database: `pkg/db/badger.go`
- Storage: `pkg/fs/local.go`, `pkg/fs/s3.go`
- Feed generation: `pkg/feed/xml.go` (RSS, filename handling)
- Filename migration: `services/migrate/migrate.go`
- Web server: `services/web/server.go`
- YouTube builder: `pkg/builder/youtube.go`
- Vimeo builder: `pkg/builder/vimeo.go`
- SoundCloud builder: `pkg/builder/soundcloud.go`
- Twitch builder: `pkg/builder/twitch.go`
- URL parsing: `pkg/builder/url.go`
- youtube-dl wrapper: `pkg/ytdl/ytdl.go`
- Hooks: `pkg/feed/hooks.go`
- API key rotation: `pkg/feed/key.go`
