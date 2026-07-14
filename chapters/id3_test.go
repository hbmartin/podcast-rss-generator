package chapters_test

import (
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
)

func fixture(name string) string {
	return filepath.Join("testdata", name)
}

func TestExtractID3Chapters(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 310, Title: "Main topic", URL: "https://example.com/topic"},
	}, chapters.ExtractID3Chapters(fixture("with_chapters.mp3")))
}

func TestExtractID3PartialCTOCKeepsUnreferencedChapters(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []chapters.Chapter{
		{Start: 310, Title: "Main topic"},
		{Start: 0, Title: "Intro"},
	}, chapters.ExtractID3Chapters(fixture("partial_ctoc.mp3")))
}

func TestExtractID3ChaptersWithoutCTOCAreSortedByStart(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 310, Title: "Main topic"},
	}, chapters.ExtractID3Chapters(fixture("no_ctoc.mp3")))
}

func TestExtractID3MissingFile(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractID3Chapters(fixture("does-not-exist.mp3")))
}

func TestExtractID3FileWithoutID3(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractID3Chapters(fixture("raw.mp3")))
}

func TestExtractID3FileWithoutCHAPFrames(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractID3Chapters(fixture("no_chap.mp3")))
}
