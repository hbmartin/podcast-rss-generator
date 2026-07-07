# Contributing

Thanks for helping improve this library. This repository is a Go module for
generating RSS 2.0 podcast feeds with Apple Podcasts tags.

## Getting Started

- Use the mise-managed toolchain from `mise.toml`.
- Keep compatibility with Go 1.24.0 and newer.
- Open pull requests against `master`.
- Update `doc.go` or Go comments when documentation needs to change; regenerate
  `README.md` with `mise run readme`.

## Required Checks

Run these before committing or opening a PR:

```sh
mise run fmt
mise run lint
mise run test
```

`mise run check` runs formatting, linting, vetting, and tests together. The CI
workflow runs the same project tasks.

## Coding Guidelines

- Follow Effective Go and the repository's existing style.
- Keep runtime dependencies minimal. Test-only dependencies are acceptable when
  they are already used in the package or materially improve test quality.
- Preserve the public API unless a breaking change is intentional and called out
  clearly in the PR.
- Use explicit error handling and wrap propagated errors with useful context.

## Testing Strategy

- Prefer examples in `examples_test.go` and `example_test.go` for positive public
  API behavior; they double as pkg.go.dev documentation.
- Put negative and edge-case public API tests in external package tests such as
  `podcast_test.go`.
- Use same-package tests such as `podcast_internals_test.go` only for internal
  behavior that cannot be reached through the public API.
- Add table-driven tests for new business logic and validation paths.
- Native fuzz coverage lives in `fuzz_test.go`; run a quick smoke check with
  `mise run fuzz-smoke`.

## Releases

Releases are published from `master` by pushing a semantic version tag such as
`v2.0.0`. Update `CHANGELOG.md` with user-visible changes before release.
