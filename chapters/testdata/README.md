# chapters test fixtures

The `*.mp3` files are 128 zero bytes followed by an ID3v2.4 tag written by
[mutagen](https://mutagen.readthedocs.io/). They mirror the fixtures the Python
`podcast-chapter-tools` test suite builds at runtime, and are used to validate
the ID3 CHAP/CTOC parser in `id3.go`.

| File | Contents |
| --- | --- |
| `with_chapters.mp3` | CTOC (ordered: chp1, chp2) + CHAP chp1 (0s, "Intro") + CHAP chp2 (310s, "Main topic", WXXX url) |
| `partial_ctoc.mp3` | CTOC referencing only chp2 + CHAP chp1 (0s) + CHAP chp2 (310s) |
| `no_ctoc.mp3` | CHAP chp2 (310s) + CHAP chp1 (0s), no CTOC |
| `no_chap.mp3` | A single TIT2 frame, no CHAP frames |
| `raw.mp3` | 128 zero bytes, no ID3 header |
