package podcast_test

import (
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
	"github.com/stretchr/testify/assert"
)

func TestItemAddSummaryTooLong(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}
	summary := ""
	for len(summary) < 4051 {
		summary += "abc ss 5 "
	}

	// act
	i.AddSummary(summary)

	// assert
	assert.Len(t, i.ISummary.Text, 4000)
}

func TestItemAddImageEmptyUrl(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	// act
	i.AddImage("")

	// assert
	assert.Nil(t, i.IImage)
}

func TestItemAddDurationZero(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}
	d := int64(0)

	// act
	i.AddDuration(d)

	// assert
	assert.EqualValues(t, "", i.IDuration)
}

func TestItemAddDurationLessThanZero(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}
	d := int64(-13)

	// act
	i.AddDuration(d)

	// assert
	assert.EqualValues(t, "", i.IDuration)
}

func TestItemAddEpisodeDefaultsToFull(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{Title: "item.title"}

	// act
	i.AddEpisode()

	// assert
	assert.NotNil(t, i.IEpisodeType)
	assert.Equal(t, "full", i.IEpisodeType.Text)
}

func TestItemAddEpisodeType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		episodeType podcast.EpisodeType
		want        string
	}{
		{"full", podcast.Full, "full"},
		{"trailer", podcast.Trailer, "trailer"},
		{"bonus", podcast.Bonus, "bonus"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			i := podcast.Item{Title: "item.title"}

			// act
			i.AddEpisodeType(tt.episodeType)

			// assert
			assert.NotNil(t, i.IEpisodeType)
			assert.Equal(t, tt.want, i.IEpisodeType.Text)
		})
	}
}

func TestItemAddDescription(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{Title: "item.title"}

	// act
	i.AddDescription("a rich <a href=\"http://example.com\">description</a>")

	// assert
	assert.EqualValues(t, "a rich <a href=\"http://example.com\">description</a>", i.Description)
}

func TestItemAddDescriptionTooLong(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{Title: "item.title"}
	d := ""
	for len(d) < 4051 {
		d += "abc ss 5 "
	}

	// act
	i.AddDescription(d)

	// assert
	assert.Len(t, string(i.Description), 4000)
}
