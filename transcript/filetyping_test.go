package transcript_test

import (
	"context"
	"path/filepath"
	"sort"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
)

func TestFileTypeString(t *testing.T) {
	t.Parallel()
	cases := map[transcript.FileType]string{
		transcript.FileTypeHTML:    "html",
		transcript.FileTypeJSON:    "json",
		transcript.FileTypeSRT:     "srt",
		transcript.FileTypeVTT:     "vtt",
		transcript.FileTypeXML:     "xml",
		transcript.FileTypeUnknown: "unknown",
	}
	for value, want := range cases {
		assert.Equal(t, want, value.String())
	}
	assert.Equal(t, "unknown", transcript.FileType(99).String())
}

func TestIdentifyFileTypeFromContent(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "no_extension_vtt")
	writeFile(t, path, "WEBVTT\n\n00:00.000 --> 00:01.000\nhi\n")
	assert.Equal(t, transcript.FileTypeVTT, transcript.IdentifyFileType(context.Background(), path))
}

func TestIdentifyFileTypesUsesContentForUnknownExtensions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	vtt := filepath.Join(dir, "transcript.txt")
	writeFile(t, vtt, "WEBVTT\n\n00:00.000 --> 00:01.000\nhi\n")
	srt := filepath.Join(dir, "episode.subtitle")
	writeFile(t, srt, "1\n00:00:00,000 --> 00:00:01,000\nhi\n\n")
	mystery := filepath.Join(dir, "mystery.bin")
	writeFile(t, mystery, "no transcript here")
	empty := filepath.Join(dir, "empty.dat")
	writeFile(t, empty, "")

	grouped := transcript.IdentifyFileTypes(context.Background(), []string{vtt, srt, mystery, empty})

	assert.Equal(t, []string{vtt}, grouped[transcript.FileTypeVTT])
	assert.Equal(t, []string{srt}, grouped[transcript.FileTypeSRT])
	unknown := grouped[transcript.FileTypeUnknown]
	sort.Strings(unknown)
	want := []string{empty, mystery}
	sort.Strings(want)
	assert.Equal(t, want, unknown)
}

func TestIdentifyFileTypesUnreadableFileIsUnknown(t *testing.T) {
	t.Parallel()
	missing := filepath.Join(t.TempDir(), "gone.txt")
	grouped := transcript.IdentifyFileTypes(context.Background(), []string{missing})
	assert.Equal(t, []string{missing}, grouped[transcript.FileTypeUnknown])
}

func TestIdentifyFileTypeUnknownExtensionUnreadable(t *testing.T) {
	t.Parallel()
	missing := filepath.Join(t.TempDir(), "gone")
	assert.Equal(t, transcript.FileTypeUnknown, transcript.IdentifyFileType(context.Background(), missing))
}
