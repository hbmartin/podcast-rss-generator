// Package podcast generates RSS 2.0 podcast feeds with common Apple Podcasts
// tags using a small Go API.
//
// Full documentation with detailed examples is located at https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2
//
// # Usage
//
// To use, `go get` and `import` the package like your typical Go library.
//
//	$ go get github.com/hbmartin/podcast-rss-generator/v2@latest
//
//	import "github.com/hbmartin/podcast-rss-generator/v2"
//
// The API exposes method receivers on feed structs that fill derived RSS and
// Apple Podcasts fields, validate required item fields, and keep generated XML
// consistent.
//
// Notably, the `Podcast.AddItem` function performs most
// of the heavy lifting by taking the `Item` input and performing
// validation, overrides and duplicate setters through the feed.
//
// Detailed examples of the API are at https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2.
//
// # Contributing
//
// See the CONTRIBUTING.md for all the details.
//
// # Go Modules
//
// This library is supported on Go 1.24.0 and higher.
//
// The repository uses Go modules and a mise-managed toolchain for local and CI
// checks.
//
// # Extensibility
//
// The exported structs remain available for callers that need direct control
// over RSS and Apple Podcasts fields. Prefer the Add methods for fields with
// formatting, validation, or derived values.
//
// # Fuzzing
//
// Native Go fuzzing covers feed encoding with XML-sensitive text, invalid UTF-8,
// long descriptions, enclosure metadata, and date bytes.
//
//	mise run fuzz-smoke
//
// To run longer fuzzing locally:
//
//	go test -run '^$' -fuzz=FuzzPodcastEncode -fuzztime=1m ./...
//
// # Roadmap
//
// Current v2 work focuses on preserving the public API, improving validation
// and safety, and keeping generated documentation and examples accurate.
//
// # Versioning
//
// Releases follow semantic versioning. Public API removals or incompatible
// behavior changes require a new major version.
//
// # Release Notes
//
// See CHANGELOG.md in the repository for release notes. The changelog is the
// source of truth for published version history.
//
// # References
//
// RSS 2.0: https://cyber.harvard.edu/rss/rss.html
//
// Podcasts: https://help.apple.com/itc/podcasts_connect/#/itca5b22233
package podcast
