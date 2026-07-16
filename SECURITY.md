# Security Policy

## Supported Versions

The latest released `v2.x` version receives security fixes.

| Version | Supported          |
| ------- | ------------------ |
| 2.x     | :white_check_mark: |
| < 2.0   | :x:                |

## Reporting a Vulnerability

Please report suspected vulnerabilities privately rather than opening a public issue.

Use GitHub's private vulnerability reporting for this repository
("Security" tab → "Report a vulnerability"), which notifies the maintainers directly.

Please include:

- A description of the issue and its impact.
- Steps to reproduce, or a minimal proof of concept.
- Any known mitigations or affected versions.

You can expect an initial acknowledgement within a few business days. Once a fix is
available, it will be released as a new tagged version and credited in the changelog
unless you prefer to remain anonymous.

Note that this is a library for generating RSS/iTunes podcast XML; it performs no
network I/O and executes no untrusted input on its own. Reports about malformed-input
handling of the feed-generation API are welcome.
