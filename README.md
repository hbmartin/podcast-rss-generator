[![Go Reference](https://pkg.go.dev/badge/github.com/hbmartin/podcast-rss-generator/v2.svg)](https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2)
[![CI](https://github.com/hbmartin/podcast-rss-generator/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/hbmartin/podcast-rss-generator/actions/workflows/ci.yml)
[![Coverage Status](https://codecov.io/gh/hbmartin/podcast-rss-generator/branch/master/graph/badge.svg)](https://codecov.io/gh/hbmartin/podcast-rss-generator)
[![Go Report Card](https://goreportcard.com/badge/github.com/hbmartin/podcast-rss-generator/v2)](https://goreportcard.com/report/github.com/hbmartin/podcast-rss-generator/v2)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

# podcast

Generate valid [RSS 2.0](https://cyber.harvard.edu/rss/rss.html) podcast feeds
with [Apple Podcasts](https://help.apple.com/itc/podcasts_connect/#/itca5b22233)
(iTunes) tags from a small, dependency-free Go API.

```go
import "github.com/hbmartin/podcast-rss-generator/v2"
```

## Why this library

- **RSS 2.0 + Apple Podcasts.** Emits both the standard RSS channel/item tags and
  the `itunes:` extension tags that Apple Podcasts requires.
- **Validation built in.** `Podcast.AddItem` enforces the required fields for
  articles and audio/video episodes and returns a typed `ItemValidationError`.
- **Derived fields set for you.** The `Add*` helpers format dates, durations,
  MIME types, and byte lengths, and fill in the duplicate iTunes tags so you
  don't have to.
- **CDATA-safe rich text.** Descriptions and summaries are rendered as CDATA, so
  HTML such as `<a href="…">` stays valid XML.
- **Flexible output.** Write a feed with `Encode(io.Writer)`, or grab it as
  `Bytes()` or `String()`.
- **Zero third-party runtime dependencies**, plus native Go fuzzing over the
  encoder.

## Contents

- [Install](#install)
- [Quick start](#quick-start)
- [Serving a feed over HTTP](#serving-a-feed-over-http)
- [Minimum requirements & validation](#minimum-requirements--validation)
- [Key types & methods](#key-types--methods)
- [Full API reference](#full-api-reference)
- [Contributing](#contributing)
- [Fuzzing](#fuzzing)
- [Versioning & releases](#versioning--releases)
- [References](#references)

## Install

```sh
go get github.com/hbmartin/podcast-rss-generator/v2@latest
```

This library supports Go 1.24.0 and higher. It uses Go modules and a
mise-managed toolchain for local and CI checks.

## Quick start

Create a podcast, add an episode, and print the feed:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hbmartin/podcast-rss-generator/v2"
)

func main() {
	now := time.Now()

	// Required channel fields: title, link, description, pubDate, lastBuildDate.
	p := podcast.New(
		"My Show",
		"https://example.com",
		"A show about interesting things",
		now, now,
	)
	p.AddAuthor("Jane Doe", "jane@example.com")
	p.AddImage("https://example.com/artwork.jpg")

	// Build an episode. AddEnclosure formats the MIME type and byte length.
	item := podcast.Item{Title: "Episode 1"}
	item.AddDescription("The very first episode.")
	item.AddPubDate(now)
	item.AddEnclosure("https://example.com/ep1.mp3", podcast.MP3, 12_345_678)
	item.AddDuration(1830) // seconds -> "30:30"

	// AddItem validates the episode and fills in derived iTunes tags.
	if _, err := p.AddItem(item); err != nil {
		fmt.Fprintln(os.Stderr, "invalid episode:", err)
		return
	}

	// Emit the RSS 2.0 feed.
	fmt.Println(p.String())
}
```

`New` defaults any zero-value time to the current UTC time and formats non-zero
values into the RSS date format. The `Add*` helpers are preferred over setting
struct fields directly because they populate the derived and duplicate fields
that Apple Podcasts expects.

## Serving a feed over HTTP

Because `Podcast.Encode` takes an `io.Writer`, you can write a feed straight to
an `http.ResponseWriter`:

```go
func feedHandler(w http.ResponseWriter, _ *http.Request) {
	p := podcast.New(
		"My Show",
		"https://example.com",
		"A show about interesting things",
		time.Now(), time.Now(),
	)
	p.AddAuthor("Jane Doe", "jane@example.com")
	p.AddAtomLink("https://example.com/feed.rss")

	item := podcast.Item{Title: "Episode 1"}
	item.AddDescription("The very first episode.")
	item.AddEnclosure("https://example.com/ep1.mp3", podcast.MP3, 12_345_678)
	if _, err := p.AddItem(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	if err := p.Encode(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

See the runnable `Example_httpHandlers` and `Example_ioWriter` examples in
[`examples_test.go`](./examples_test.go) for full, verified output.

## Minimum requirements & validation

`Podcast.AddItem` rejects episodes that are missing required fields and returns
an `*ItemValidationError` (which unwraps to one of the sentinel errors below).

| Episode kind | Required fields |
| --- | --- |
| Article | `Title`, `Description`, `Link` |
| Audio / video / download | `Title`, `Description`, `Enclosure` (URL, `Type`, and `Length` all required) |

Some fields are **always overwritten** by `AddItem` — don't set them yourself:

- `GUID`
- `PubDateFormatted`
- `AuthorFormatted`
- `Enclosure.TypeFormatted`
- `Enclosure.LengthFormatted`

Sentinel errors returned by feed and item validation:

```go
var (
	ErrPodcastRequired          = errors.New("podcast is required")
	ErrWriterRequired           = errors.New("writer is required")
	ErrTitleDescriptionRequired = errors.New("title and description are required")
	ErrEnclosureURLRequired     = errors.New("enclosure url is required")
	ErrEnclosureTypeRequired    = errors.New("enclosure type is required")
	ErrLinkRequired             = errors.New("link is required when not using enclosure")
)
```

You can test for a specific cause with `errors.Is`:

```go
if _, err := p.AddItem(item); err != nil {
	if errors.Is(err, podcast.ErrEnclosureURLRequired) {
		// handle the missing enclosure URL
	}
}
```

## Key types & methods

The full, always-current API is on
[pkg.go.dev](https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2).
The most commonly used pieces:

| Type | Purpose |
| --- | --- |
| `Podcast` | The RSS channel. Create with `New(...)`. |
| `Item` | A single episode. |
| `Enclosure` / `EnclosureType` | The downloadable media file and its MIME type (`M4A`, `M4V`, `MP4`, `MP3`, `MOV`, `PDF`, `EPUB`). |
| `PodcastType` | `Episodic` (default) or `Serial`. |
| `EpisodeType` | `Full` (default), `Trailer`, or `Bonus`. |
| `ItemValidationError` | Typed error explaining why an item was rejected. |

Frequently used methods:

| On `Podcast` | On `Item` |
| --- | --- |
| `New(title, link, description, pubDate, lastBuildDate)` | `AddDescription(text)` |
| `AddAuthor(name, email)` | `AddSummary(text)` |
| `AddImage(url)` | `AddPubDate(t)` |
| `AddCategory(category, subCategories)` | `AddEnclosure(url, type, length)` |
| `AddSummary(text)` / `AddSubTitle(text)` | `AddDuration(seconds)` |
| `AddType(podcastType)` | `AddImage(url)` |
| `AddAtomLink(href)` | `AddEpisode()` |
| `AddItem(item) (int, error)` | `AddEpisodeType(episodeType)` |
| `Encode(w)` / `Bytes()` / `String()` | |

The exported structs remain available for callers that need direct control over
RSS and Apple Podcasts fields. Prefer the `Add*` methods for any field with
formatting, validation, or derived values.

## Full API reference

Complete documentation, every type, and runnable examples are generated from
the source and hosted on pkg.go.dev:

> **<https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2>**

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details. Before opening a PR, run:

```sh
mise run fmt   # format
mise run lint  # linters
mise run test  # unit tests
```

## Fuzzing

Native Go fuzzing covers feed encoding with XML-sensitive text, invalid UTF-8,
long descriptions, enclosure metadata, and date bytes.

```sh
mise run fuzz-smoke
```

To run longer fuzzing locally:

```sh
go test -run '^$' -fuzz=FuzzPodcastEncode -fuzztime=1m ./...
```

## Versioning & releases

Releases follow [semantic versioning](https://semver.org/). Public API removals
or incompatible behavior changes require a new major version.
[CHANGELOG.md](CHANGELOG.md) is the source of truth for published version
history.

## References

- RSS 2.0 specification: <https://cyber.harvard.edu/rss/rss.html>
- Apple Podcasts requirements: <https://help.apple.com/itc/podcasts_connect/#/itca5b22233>
- Apple Podcasts category list: <https://help.apple.com/itc/podcasts_connect/#/itc9267a2f12>
