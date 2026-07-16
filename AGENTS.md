## Development Guidelines

This repository is a single-package Go **library** (`package podcast`) that generates
fully compliant iTunes and RSS 2.0 podcast feeds. There is no `main` package, no
binary, and no network/IO layer — all code lives in `.go` files at the repository root.

Important: Write tests for any changes you make, and fill in any test coverage gaps.

## Go Coding Standards & Idioms
When writing or refactoring Go code in this repository, strictly adhere to the following rules:

1. **Error Handling:**
   - Never use `panic()` unless explicitly initializing a package-level variable that must not fail (e.g., compiling a regex).
   - Always check errors explicitly (`if err != nil`).
   - Wrap errors with context using `fmt.Errorf("doing something: %w", err)` to maintain the stack trace.

2. **Concurrency (when applicable):**
   - This library is synchronous and has no goroutines today. If you add concurrency, pass `context.Context` as the first argument to functions that do I/O or heavy processing.
   - Guard goroutines with a `sync.WaitGroup` or an `errgroup.Group`. Do not leave orphaned goroutines.
   - Ensure channels are closed by the sender, not the receiver.

3. **Style & Formatting:**
   - Follow [Effective Go](https://go.dev/doc/effective_go).
   - Use short variable names for limited scopes (e.g., `r` for `http.Request`).
   - Keep interfaces small, typically 1-3 methods, defined where they are consumed.
   - Avoid package-level state/globals.

4. **Testing:**
   - Write table-driven tests for all business logic.
   - Use the standard `testing` package. This package's existing tests use `testify` and godoc-style `Example` functions — match the style of the file you are editing.
   - Name test files `[filename]_test.go`.
   - Testable `Example` functions double as documentation on pkg.go.dev — keep their `// Output:` blocks accurate.

### Code Quality
- Write clean, idiomatic Go code following Go conventions and best practices.
- Ensure proper error handling and meaningful error messages.
- Follow the existing code style and patterns in the repository.

### Testing and Quality Assurance
- **CRITICAL**: Always run ALL of the following before making a commit or opening a PR:
  1. `go fmt ./...` — Format all Go files (or `make fmt`)
  2. `golangci-lint run` — Run all configured linters and formatters (or `make lint`)
  3. `go test -race ./...` — Run all unit tests (or `make test`)
- Ensure ALL tests pass AND ALL linting checks pass before committing.
- If you change dependencies or the Go version, run `make tidy` (`go mod tidy && go mod vendor`) so `vendor/` and `go.mod`/`go.sum` stay in sync — CI enforces this.

## Formatting and Linting Requirements

This project uses golangci-lint (v2 config in `.golangci.yml`) with `gofmt` + `goimports`
formatters. Common requirements: proper spacing around operators, correct struct field
alignment, import ordering (standard library, third-party, local), no trailing whitespace,
and a space after commas.

**Always run `make fmt`, `make lint`, AND `make test` after making ANY code changes.**

## Key File References

Flat single-package library at the repository root (`package podcast`):

- Core feed type + `New`/`AddItem`/`Bytes`/`Encode`/`String`: `podcast.go`
- Episode type + `AddEnclosure`/`AddImage`/`AddDuration`/`AddSummary`: `item.go`
- Media enclosure element: `enclosure.go`
- iTunes-specific elements: `itunes.go`
- Podcast type/category enum: `type.go`
- Sub-elements: `author.go`, `atomlink.go`, `description.go`, `image.go`, `textinput.go`
- Package documentation (rendered on pkg.go.dev; README is generated from this): `doc.go`
- Fuzz harnesses for exported funcs: `fuzz.go`
- Tests + godoc examples: `*_test.go`, `example_test.go`, `examples_test.go`

## Build, Test & Release

- Build/vet: `make build` / `go vet ./...`
- Test + coverage: `make test`, `make cover`
- Lint: `make lint`
- CI runs on pushes/PRs to `master` (`.github/workflows/ci.yml`).
- Release = push a semantic-version tag (`vMAJOR.MINOR.PATCH`) from `master`. This is a
  v2 module, so the module path carries a `/v2` suffix. The `go-mod-publish` workflow
  warms the Go module proxy on tag push so the release appears on pkg.go.dev.
