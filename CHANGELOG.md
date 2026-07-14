# Changelog

All notable changes to this project are documented here.

## Unreleased

- Add `transcript` subpackage: convert podcast transcripts from SRT, WebVTT,
  Podlove simple-transcripts XML, HTML, and PodcastIndex JSON into
  PodcastIndex transcript JSON, individually or in bulk. Ported from
  [podcast-transcript-convert](https://github.com/hbmartin/podcast-transcript-convert).
- Add `transcript/overcastdb` subpackage: read transcript paths and episode
  metadata from an overcast-to-sqlite database (isolated so the pure-Go
  SQLite driver is only pulled in by consumers that use it).
- Add `chapters` subpackage: extract and transform podcast chapters between
  PodcastIndex (PCI) chapters, Podlove Simple Chapters (PSC), description
  timestamp embeds, and ID3v2 CHAP frames. Ported from
  [podcast-chapter-tools](https://github.com/hbmartin/podcast-chapter-tools).
- Add partial Podcasting 2.0 (podcast namespace) support: channel-level
  `podcast:guid`, `podcast:medium`, `podcast:locked`, and `podcast:person`;
  item-level `podcast:transcript`, `podcast:chapters`, `podcast:person`, and
  `podcast:socialInteract`. The `xmlns:podcast` declaration is emitted only
  when a podcast-namespace tag is set. Includes a dependency-free
  `NewFeedGUID` helper implementing the spec's UUIDv5 derivation.
- Add inspectable validation errors for podcast and item failures.
- Make podcast encoding safe for nil receivers, nil writers, and zero-value
  podcast values.
- Clone nested item pointer fields when adding items to avoid caller-owned
  state aliasing.
- Replace legacy go-fuzz targets with native Go fuzzing.
- Remove vendored test dependencies and refresh the test dependency graph.
- Expand lint coverage and refresh contributor, README, and agent-facing docs.

## v2.0.0

- Change module path to `github.com/hbmartin/podcast-rss-generator/v2` for semantic import versioning.
- Publish as a Go library module with CI and tag workflows focused on build, vet, lint, and tests.
- Raise the minimum supported Go version to 1.24.0.

## v1.4.2

- Slim down Go Modules for consumers (#32).

## v1.4.1

- Implement fuzz logic testing of exported funcs (#31).
- Upgrade CICD Pipeline Tooling (#31).
- Update documentation for 1.x and 2.3 (#31).
- Allow godoc2ghmd to run without network (#31).

## v1.4.0

- Add Go Modules, Update vendor folder (#26, #25).
- Add C.I. GitHub Actions (#25).
- Add additional error checks found by linters (#25).
- Go Fmt enclosure_test.go (#25).

## v1.3.2

- Correct count len of UTF8 strings (#9).
- Implement duration parser (#8).
- Fix Github and GoDocs Markdown (#14).
- Move podcast.go Private Methods to Respected Files (#12).
- Allow providing GUID on Podcast (#15).

## v1.3.1

- Increased iTunes compliance after feedback from Apple.
- Specified what categories should be set with AddCategory().
- Enforced title and link as part of Image.
- Added Podcast.AddAtomLink() for more broad compliance to readers.

## v1.3.0

- Fix Item.Duration being set incorrectly.
- Change Item.AddEnclosure() parameter definition (Bytes not Seconds).
- Add Item.AddDuration formatting and override.
- Add more documentation surrounding Item.Enclosure{}.

## v1.2.1

- Add Podcast.AddSubTitle() and truncating to 64 chars.
- Add guards to protect against empty fields.

## v1.2.0

- Add Podcast.AddPubDate() and Podcast.AddLastBuildDate() overrides.
- Add Item.AddImage() to mask some cumbersome addition of IImage.
- Add Item.AddPubDate to simplify datetime setters.
- Add more examples, mostly around Item struct.
- Tweak documentation.

## v1.1.0

- Enable CDATA in ISummary fields for Podcast and Channel.

## v1.0.0

- Initial release.
- Full documentation, full examples, and complete code coverage.
