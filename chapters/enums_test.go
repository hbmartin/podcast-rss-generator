package chapters_test

import (
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
)

func TestChapterTypeString(t *testing.T) {
	t.Parallel()
	cases := map[chapters.ChapterType]string{
		chapters.ChapterTypeAI:          "ai",
		chapters.ChapterTypeDescription: "description",
		chapters.ChapterTypeGemini:      "gemini",
		chapters.ChapterTypeID3:         "id3",
		chapters.ChapterTypeOpenAI:      "openai",
		chapters.ChapterTypePCI:         "pci",
		chapters.ChapterTypePSC:         "psc",
		chapters.ChapterTypeSonnet:      "sonnet",
	}
	for value, want := range cases {
		assert.Equal(t, want, value.String())
	}
	assert.Empty(t, chapters.ChapterType(99).String())
}
