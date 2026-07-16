// Package podcast generates a fully compliant iTunes and RSS 2.0 podcast feed
// for GoLang using a simple API.
//
// Full documentation with detailed examples located at https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2
//
// # Usage
//
// To use, `go get` and `import` the package like your typical GoLang library.
//
//	$ go get github.com/hbmartin/podcast-rss-generator/v2
//
//	import "github.com/hbmartin/podcast-rss-generator/v2"
//
// The API exposes a number of method receivers on structs that implements the
// logic required to comply with the specifications and ensure a compliant feed.
// A number of overrides occur to help with iTunes visibility of your episodes.
//
// Notably, the `Podcast.AddItem` function performs most
// of the heavy lifting by taking the `Item` input and performing
// validation, overrides and duplicate setters through the feed.
//
// Full detailed Examples of the API are at https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2.
//
// # Contributing
//
// See the CONTRIBUTING.md for all the details.
//
// # Go Modules
//
// This library requires Go 1.21 or higher and is a Go module. Import it at the
// /v2 module path shown above (semantic import versioning). Dependencies are also
// vendored in the vendor/ folder.
//
// If you hit a problem on a supported runtime, please open an Issue.
//
// # Extensibility
//
// You are not restricted in having full control over your feeds. You may choose to
// skip the API methods and instead use the structs directly.  The fields have been
// grouped by RSS 2.0 and iTunes fields with iTunes specific fields all prefixed with
// the letter `I`.
//
// # Fuzzing Inputs
//
// `go-fuzz` has been added in 1.4.1, covering all exported API methods.  They have been
// ran extensively and no issues have come out of them yet (most tests were ran overnight,
// over about 11 hours with zero crashes).
//
// If you wish to help fuzz the inputs, with Go 1.13 or later you can run `go-fuzz` on any
// of the inputs.
//
//	go get -u github.com/dvyukov/go-fuzz/go-fuzz
//	go get -u github.com/dvyukov/go-fuzz/go-fuzz-build
//	go get github.com/hbmartin/podcast-rss-generator/v2
//	cd $GOPATH/src/github.com/hbmartin/podcast-rss-generator
//	go-fuzz-build
//	go-fuzz -func FuzzPodcastAddItem
//
// To obtain a list of available funcs to pass, just run `go-fuzz` without any parameters:
//
//	$ go-fuzz
//	2020/02/13 07:27:32 -func flag not provided, but multiple fuzz functions available:
//	FuzzItemAddDuration, FuzzItemAddEnclosure, FuzzItemAddImage, FuzzItemAddPubDate,
//	FuzzItemAddSummary, FuzzPodcastAddAtomLink, FuzzPodcastAddAuthor, FuzzPodcastAddCategory,
//	FuzzPodcastAddImage, FuzzPodcastAddItem, FuzzPodcastAddLastBuildDate, FuzzPodcastAddPubDate,
//	FuzzPodcastAddSubTitle, FuzzPodcastAddSummary, FuzzPodcastBytes, FuzzPodcastEncode,
//	FuzzPodcastNew
//
// If you do find an issue, please raise an issue immediately and I will quickly address.
//
// # Roadmap
//
// This is a maintained fork of the original github.com/eduncan911/podcast
// library. The core feed-generation API is stable and considered production-ready.
// Contributions are welcome via PRs.
//
// # Versioning
//
// We follow SemVer. This library is published as a v2 Go module (the /v2 import
// path), which reflects the new module location under this fork. The public API is
// compatible with the original 1.x releases, so existing call sites work unchanged
// after updating the import path.
//
// # Release Notes
//
// v2.0.0
//   - Republish as a maintained fork at github.com/hbmartin/podcast-rss-generator (/v2 module path).
//   - Modernize CI: build/test/vet/lint on Go 1.21+ across Linux, macOS and Windows.
//   - Remove unrelated application build/release tooling; publish as a Go module via tags.
//   - Public feed-generation API unchanged from 1.x.
//
// v1.4.2
//   - Slim down Go Modules for consumers (#32)
//
// v1.4.1
//   - Implement fuzz logic testing of exported funcs (#31)
//   - Upgrade CICD Pipeline Tooling (#31)
//   - Update documentation for 1.x and 2.3 (#31)
//   - Allow godoc2ghmd to run without network (#31)
//
// v1.4.0
//   - Add Go Modules, Update vendor folder (#26, #25)
//   - Add C.I. GitHub Actions (#25)
//   - Add additional error checks found by linters (#25)
//   - Go Fmt enclosure_test.go (#25)
//
// v1.3.2
//   - Correct count len of UTF8 strings (#9)
//   - Implement duration parser (#8)
//   - Fix Github and GoDocs Markdown (#14)
//   - Move podcast.go Private Methods to Respected Files (#12)
//   - Allow providing GUID on Podcast (#15)
//
// v1.3.1
//   - increased itunes compliance after feedback from Apple:
//   - specified what categories should be set with AddCategory().
//   - enforced title and link as part of Image.
//   - added Podcast.AddAtomLink() for more broad compliance to readers.
//
// v1.3.0
//   - fixes Item.Duration being set incorrectly.
//   - changed Item.AddEnclosure() parameter definition (Bytes not Seconds!).
//   - added Item.AddDuration formatting and override.
//   - added more documentation surrounding Item.Enclosure{}
//
// v1.2.1
//   - added Podcast.AddSubTitle() and truncating to 64 chars.
//   - added a number of Guards to protect against empty fields.
//
// v1.2.0
//   - added Podcast.AddPubDate() and Podcast.AddLastBuildDate() overrides.
//   - added Item.AddImage() to mask some cumbersome addition of IImage.
//   - added Item.AddPubDate to simply datetime setters.
//   - added more examples (mostly around Item struct).
//   - tweaked some documentation.
//
// v1.1.0
//   - Enabling CDATA in ISummary fields for Podcast and Channel.
//
// v1.0.0
//   - Initial release.
//   - Full documentation, full examples and complete code coverage.
//
// # References
//
// RSS 2.0: https://cyber.harvard.edu/rss/rss.html
//
// Podcasts: https://help.apple.com/itc/podcasts_connect/#/itca5b22233
package podcast
