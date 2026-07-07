# Changelog

All notable changes to this project are documented here.

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
