# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0]

Republished as a maintained fork of `github.com/eduncan911/podcast`.

### Changed
- Module path is now `github.com/hbmartin/podcast-rss-generator/v2` (semantic import
  versioning). Update imports to the `/v2` path; the public API is unchanged from 1.x.
- Minimum Go version raised to 1.21.
- Modernized CI: `build` / `vet` / `test -race` / `golangci-lint` run on pushes and
  pull requests to `master` across Linux, macOS, and Windows.

### Removed
- Removed unrelated application build/release tooling (goreleaser config and the
  Docker image release workflow) that had been copied in from another project. This
  is a library and is released as a Go module via version tags.

## [1.4.2]
- Slim down Go Modules for consumers (#32)

## [1.4.1]
- Implement fuzz logic testing of exported funcs (#31)
- Upgrade CI/CD pipeline tooling (#31)
- Update documentation (#31)
- Allow godoc2ghmd to run without network (#31)

## [1.4.0]
- Add Go Modules, update vendor folder (#26, #25)
- Add CI GitHub Actions (#25)
- Add additional error checks found by linters (#25)

## [1.3.2]
- Correct count len of UTF8 strings (#9)
- Implement duration parser (#8)
- Fix GitHub and GoDocs Markdown (#14)
- Move podcast.go private methods to respective files (#12)
- Allow providing GUID on Podcast (#15)

## [1.3.1]
- Increased iTunes compliance after feedback from Apple (categories, image title/link)
- Added `Podcast.AddAtomLink()` for broader reader compliance

## [1.3.0]
- Fix `Item.Duration` being set incorrectly
- Change `Item.AddEnclosure()` parameter definition (Bytes not Seconds)
- Add `Item.AddDuration` formatting and override

## [1.2.1]
- Add `Podcast.AddSubTitle()` and truncate to 64 chars
- Add guards to protect against empty fields

## [1.2.0]
- Add `Podcast.AddPubDate()` and `Podcast.AddLastBuildDate()` overrides
- Add `Item.AddImage()` and `Item.AddPubDate()`

## [1.1.0]
- Enable CDATA in ISummary fields for Podcast and Channel

## [1.0.0]
- Initial release with full documentation, examples, and complete code coverage

[2.0.0]: https://github.com/hbmartin/podcast-rss-generator/releases/tag/v2.0.0
