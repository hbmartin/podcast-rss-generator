package podcast_test

import (
	"strings"
	"testing"
	"unicode/utf8"

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

func TestItemAddDescriptionLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		in        string
		wantRunes int
	}{
		{
			name:      "keeps descriptions over old limit",
			in:        strings.Repeat("a", 4001),
			wantRunes: 4001,
		},
		{
			name:      "keeps description at new limit",
			in:        strings.Repeat("a", 10000),
			wantRunes: 10000,
		},
		{
			name:      "truncates description over new limit",
			in:        strings.Repeat("é", 10001),
			wantRunes: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := podcast.Item{}
			i.AddDescription(tt.in)

			got := string(i.Description)
			if count := utf8.RuneCountInString(got); count != tt.wantRunes {
				t.Fatalf("AddDescription() rune count = %d, want %d", count, tt.wantRunes)
			}
			if !utf8.ValidString(got) {
				t.Fatal("AddDescription() returned invalid UTF-8")
			}
		})
	}
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
	assert.Empty(t, i.IDuration)
}

func TestItemAddEpisodeType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		episodeType podcast.EpisodeType
		want        string
	}{
		{name: "full", episodeType: podcast.Full, want: "full"},
		{name: "trailer", episodeType: podcast.Trailer, want: "trailer"},
		{name: "bonus", episodeType: podcast.Bonus, want: "bonus"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := podcast.Item{
				Title:       "item.title",
				Description: "item.desc",
				Link:        "http://example.com/article.html",
			}

			i.AddEpisodeType(tt.episodeType)

			if assert.NotNil(t, i.IEpisodeType) {
				assert.Equal(t, tt.want, i.IEpisodeType.Text)
			}
		})
	}
}

func TestItemAddEpisode(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	// act
	i.AddEpisode()

	// assert
	if assert.NotNil(t, i.IEpisodeType) {
		assert.Equal(t, "full", i.IEpisodeType.Text)
	}
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
	assert.Empty(t, i.IDuration)
}
