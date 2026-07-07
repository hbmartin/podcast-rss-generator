package podcast_test

import (
	"strings"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const podcastNSAttr = `xmlns:podcast="https://podcastindex.org/namespace/1.0"`

func TestPodcastNamespaceDeclaredOnlyWhenUsed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(p *podcast.Podcast, i *podcast.Item)
		want   bool
	}{
		{
			name:   "no podcast tags",
			mutate: func(_ *podcast.Podcast, _ *podcast.Item) {},
			want:   false,
		},
		{
			name: "channel guid",
			mutate: func(p *podcast.Podcast, _ *podcast.Item) {
				p.SetPodcastGUID("917393e3-1b1e-5cef-ace4-edaa54e1f810")
			},
			want: true,
		},
		{
			name: "channel medium",
			mutate: func(p *podcast.Podcast, _ *podcast.Item) {
				p.SetMedium(podcast.MediumMusic)
			},
			want: true,
		},
		{
			name: "channel locked",
			mutate: func(p *podcast.Podcast, _ *podcast.Item) {
				p.SetLocked(true, "")
			},
			want: true,
		},
		{
			name: "channel person",
			mutate: func(p *podcast.Podcast, _ *podcast.Item) {
				p.AddPerson("Jane Doe", "", "", "", "")
			},
			want: true,
		},
		{
			name: "item transcript",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddTranscript("https://example.com/1.vtt", "text/vtt", "", "")
			},
			want: true,
		},
		{
			name: "item chapters",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddChapters("https://example.com/1.json", "application/json+chapters")
			},
			want: true,
		},
		{
			name: "item person",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddPerson("Jane Doe", "guest", "", "", "")
			},
			want: true,
		},
		{
			name: "item social interact",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddSocialInteract("https://example.social/@dave/1", "activitypub", "@dave")
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := podcast.New("t", "l", "d", zeroDate, zeroDate)
			i := newValidItem()
			tt.mutate(&p, &i)
			_, err := p.AddItem(i)
			require.NoError(t, err)

			got := p.String()
			if tt.want {
				assert.Contains(t, got, podcastNSAttr)
			} else {
				assert.NotContains(t, got, "xmlns:podcast")
				// The legacy root element is untouched.
				assert.Contains(
					t, got,
					`<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">`,
				)
			}
		})
	}
}

func TestChannelPodcastTagsRender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(p *podcast.Podcast)
		want   string
	}{
		{
			name: "guid",
			mutate: func(p *podcast.Podcast) {
				p.SetPodcastGUID("917393e3-1b1e-5cef-ace4-edaa54e1f810")
			},
			want: "<podcast:guid>917393e3-1b1e-5cef-ace4-edaa54e1f810</podcast:guid>",
		},
		{
			name:   "medium",
			mutate: func(p *podcast.Podcast) { p.SetMedium(podcast.MediumMusic) },
			want:   "<podcast:medium>music</podcast:medium>",
		},
		{
			name:   "locked yes with owner",
			mutate: func(p *podcast.Podcast) { p.SetLocked(true, "jane@example.com") },
			want:   `<podcast:locked owner="jane@example.com">yes</podcast:locked>`,
		},
		{
			name:   "locked no without owner",
			mutate: func(p *podcast.Podcast) { p.SetLocked(false, "") },
			want:   "<podcast:locked>no</podcast:locked>",
		},
		{
			name: "person with all attributes",
			mutate: func(p *podcast.Podcast) {
				p.AddPerson(
					"Jane Doe", "guest", "cast",
					"https://example.com/jane.jpg", "https://example.com/jane",
				)
			},
			want: `<podcast:person role="guest" group="cast" ` +
				`img="https://example.com/jane.jpg" href="https://example.com/jane">` +
				`Jane Doe</podcast:person>`,
		},
		{
			name:   "person with defaults omitted",
			mutate: func(p *podcast.Podcast) { p.AddPerson("Jane Doe", "", "", "", "") },
			want:   "<podcast:person>Jane Doe</podcast:person>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := podcast.New("t", "l", "d", zeroDate, zeroDate)
			tt.mutate(&p)
			assert.Contains(t, p.String(), tt.want)
		})
	}
}

func TestItemPodcastTagsRender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(i *podcast.Item)
		want   string
	}{
		{
			name: "transcript with all attributes",
			mutate: func(i *podcast.Item) {
				i.AddTranscript("https://example.com/1.vtt", "text/vtt", "en", "captions")
			},
			want: `<podcast:transcript url="https://example.com/1.vtt" type="text/vtt" ` +
				`language="en" rel="captions"></podcast:transcript>`,
		},
		{
			name: "transcript minimal",
			mutate: func(i *podcast.Item) {
				i.AddTranscript("https://example.com/1.srt", "application/x-subrip", "", "")
			},
			want: `<podcast:transcript url="https://example.com/1.srt" ` +
				`type="application/x-subrip"></podcast:transcript>`,
		},
		{
			name: "chapters",
			mutate: func(i *podcast.Item) {
				i.AddChapters("https://example.com/1.json", "application/json+chapters")
			},
			want: `<podcast:chapters url="https://example.com/1.json" ` +
				`type="application/json+chapters"></podcast:chapters>`,
		},
		{
			name: "person",
			mutate: func(i *podcast.Item) {
				i.AddPerson("John Smith", "guest", "", "", "https://example.com/john")
			},
			want: `<podcast:person role="guest" href="https://example.com/john">` +
				`John Smith</podcast:person>`,
		},
		{
			name: "social interact",
			mutate: func(i *podcast.Item) {
				i.AddSocialInteract("https://example.social/@dave/1", "activitypub", "@dave")
			},
			want: `<podcast:socialInteract uri="https://example.social/@dave/1" ` +
				`protocol="activitypub" accountId="@dave"></podcast:socialInteract>`,
		},
		{
			name: "social interact without account id",
			mutate: func(i *podcast.Item) {
				i.AddSocialInteract("https://example.social/@dave/1", "activitypub", "")
			},
			want: `<podcast:socialInteract uri="https://example.social/@dave/1" ` +
				`protocol="activitypub"></podcast:socialInteract>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := podcast.New("t", "l", "d", zeroDate, zeroDate)
			i := newValidItem()
			tt.mutate(&i)
			_, err := p.AddItem(i)
			require.NoError(t, err)
			assert.Contains(t, p.String(), tt.want)
		})
	}
}

func TestPodcastSettersNoOpOnMissingRequiredInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(p *podcast.Podcast, i *podcast.Item)
	}{
		{
			name:   "empty guid",
			mutate: func(p *podcast.Podcast, _ *podcast.Item) { p.SetPodcastGUID("") },
		},
		{
			name:   "blank channel person name",
			mutate: func(p *podcast.Podcast, _ *podcast.Item) { p.AddPerson("  ", "host", "", "", "") },
		},
		{
			name:   "blank item person name",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) { i.AddPerson("", "", "", "", "") },
		},
		{
			name: "transcript without url",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddTranscript("", "text/vtt", "", "")
			},
		},
		{
			name: "transcript without type",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddTranscript("https://example.com/1.vtt", "", "", "")
			},
		},
		{
			name: "chapters without url",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddChapters("", "application/json+chapters")
			},
		},
		{
			name: "chapters without type",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddChapters("https://example.com/1.json", "")
			},
		},
		{
			name: "social interact without uri",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddSocialInteract("", "activitypub", "@dave")
			},
		},
		{
			name: "social interact without protocol",
			mutate: func(_ *podcast.Podcast, i *podcast.Item) {
				i.AddSocialInteract("https://example.social/@dave/1", "", "@dave")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := podcast.New("t", "l", "d", zeroDate, zeroDate)
			i := newValidItem()
			tt.mutate(&p, &i)
			_, err := p.AddItem(i)
			require.NoError(t, err)

			got := p.String()
			assert.NotContains(t, got, "xmlns:podcast")
			assert.NotContains(t, got, "<podcast:")
		})
	}
}

func TestAddPersonTruncatesLongNames(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)
	longName := strings.Repeat("n", 129)
	p.AddPerson(longName, "", "", "", "")

	require.Len(t, p.PPersons, 1)
	assert.Equal(t, strings.Repeat("n", 128), p.PPersons[0].Name)
}

func TestMediumString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		medium podcast.Medium
		want   string
	}{
		{podcast.MediumPodcast, "podcast"},
		{podcast.MediumMusic, "music"},
		{podcast.MediumVideo, "video"},
		{podcast.MediumFilm, "film"},
		{podcast.MediumAudiobook, "audiobook"},
		{podcast.MediumNewsletter, "newsletter"},
		{podcast.MediumBlog, "blog"},
		{podcast.MediumPublisher, "publisher"},
		{podcast.MediumCourse, "course"},
		{podcast.MediumMixed, "mixed"},
		{podcast.MediumPodcastList, "podcastL"},
		{podcast.MediumMusicList, "musicL"},
		{podcast.MediumVideoList, "videoL"},
		{podcast.MediumFilmList, "filmL"},
		{podcast.MediumAudiobookList, "audiobookL"},
		{podcast.MediumNewsletterList, "newsletterL"},
		{podcast.MediumBlogList, "blogL"},
		{podcast.Medium(999), "podcast"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.medium.String())
		})
	}
}

func TestNewFeedGUID(t *testing.T) {
	t.Parallel()

	// Reference value from the podcast-namespace guid.md specification.
	const specGUID = "917393e3-1b1e-5cef-ace4-edaa54e1f810"

	tests := []struct {
		name    string
		feedURL string
		want    string
	}{
		{"spec vector with scheme", "https://mp3s.nashownotes.com/pc20rss.xml", specGUID},
		{"http scheme stripped", "http://mp3s.nashownotes.com/pc20rss.xml", specGUID},
		{"no scheme", "mp3s.nashownotes.com/pc20rss.xml", specGUID},
		{"trailing slashes stripped", "https://mp3s.nashownotes.com/pc20rss.xml//", specGUID},
		{"empty input", "", ""},
		{"scheme only", "https://", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, podcast.NewFeedGUID(tt.feedURL))
		})
	}
}

func TestAddItemSnapshotsPodcastTags(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)

	item := newValidItem()
	item.AddTranscript("https://example.com/1.vtt", "text/vtt", "en", "captions")
	item.AddChapters("https://example.com/1.json", "application/json+chapters")
	item.AddPerson("Jane Doe", "guest", "", "", "")
	item.AddSocialInteract("https://example.social/@dave/1", "activitypub", "@dave")

	_, err := p.AddItem(item)
	require.NoError(t, err)

	// Mutating the caller's structs after AddItem must not change the feed.
	item.PTranscripts[0].URL = "https://mutated.example.com"
	item.PChapters.URL = "https://mutated.example.com"
	item.PPersons[0].Name = "Mutated Name"
	item.PSocialInteracts[0].URI = "https://mutated.example.com"

	got := p.String()
	assert.NotContains(t, got, "mutated")
	assert.NotContains(t, got, "Mutated")
	assert.Contains(t, got, `url="https://example.com/1.vtt"`)
	assert.Contains(t, got, `url="https://example.com/1.json"`)
	assert.Contains(t, got, ">Jane Doe</podcast:person>")
	assert.Contains(t, got, `uri="https://example.social/@dave/1"`)
}
