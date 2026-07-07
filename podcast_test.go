package podcast_test

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hbmartin/podcast-rss-generator/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	createdDate = time.Date(2017, time.February, 1, 8, 21, 52, 0, time.UTC)
	updatedDate = createdDate.AddDate(0, 0, 5)
	pubDate     = createdDate.AddDate(0, 0, 3)
	zeroDate    = time.Time{}
)

func TestNewNonZeroDates(t *testing.T) {
	t.Parallel()

	// arrange
	ti, l, d := "title", "link", "description"

	// act
	p := podcast.New(ti, l, d, createdDate, updatedDate)

	// assert
	assert.Equal(t, ti, p.Title)
	assert.Equal(t, l, p.Link)
	assert.Equal(t, podcast.Description(d), p.Description)
	assert.GreaterOrEqual(t, createdDate.Format(time.RFC1123Z), p.PubDate)
	assert.GreaterOrEqual(t, updatedDate.Format(time.RFC1123Z), p.LastBuildDate)
}

func TestNewZeroDates(t *testing.T) {
	t.Parallel()

	// arrange
	ti, l, d := "title", "link", "description"

	// act
	p := podcast.New(ti, l, d, zeroDate, zeroDate)

	// assert
	now := time.Now().UTC().Format(time.RFC1123Z)
	assert.Equal(t, ti, p.Title)
	assert.Equal(t, l, p.Link)
	assert.Equal(t, podcast.Description(d), p.Description)
	// ensure time.Now().UTC() is set, or close to it
	assert.GreaterOrEqual(t, now, p.PubDate)
	assert.GreaterOrEqual(t, now, p.LastBuildDate)
}

func TestAddAuthor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		authorName         string
		email              string
		wantManagingEditor string
		wantOwner          bool
	}{
		{
			name:               "name and email set owner",
			authorName:         "the name",
			email:              "me@test.com",
			wantManagingEditor: "me@test.com (the name)",
			wantOwner:          true,
		},
		{
			name: "empty email ignores author",
		},
		{
			name:               "empty name skips incomplete owner",
			email:              "me@test.com",
			wantManagingEditor: "me@test.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := podcast.New("title", "link", "description", zeroDate, zeroDate)
			p.AddAuthor(tt.authorName, tt.email)

			assert.Equal(t, tt.wantManagingEditor, p.ManagingEditor)
			assert.Equal(t, tt.wantManagingEditor, p.IAuthor)
			if !tt.wantOwner {
				assert.Nil(t, p.IOwner)
				return
			}
			if assert.NotNil(t, p.IOwner) {
				assert.Equal(t, tt.authorName, p.IOwner.Name)
				assert.Equal(t, tt.email, p.IOwner.Email)
			}
		})
	}
}

func TestAddAuthorEncodesITunesOwner(t *testing.T) {
	t.Parallel()

	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	p.AddAuthor("the name", "me@test.com")

	got := p.String()
	assert.Contains(t, got, "<itunes:owner>")
	assert.Contains(t, got, "<itunes:name>the name</itunes:name>")
	assert.Contains(t, got, "<itunes:email>me@test.com</itunes:email>")
	assert.Contains(t, got, "</itunes:owner>")
}

func TestAddAuthorClearsITunesOwnerWhenNameIsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		initialName        string
		initialEmail       string
		nextName           string
		nextEmail          string
		wantManagingEditor string
	}{
		{
			name:               "replaces complete owner with email-only author",
			initialName:        "the name",
			initialEmail:       "old@test.com",
			nextEmail:          "new@test.com",
			wantManagingEditor: "new@test.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := podcast.New("title", "link", "description", zeroDate, zeroDate)
			p.AddAuthor(tt.initialName, tt.initialEmail)
			p.AddAuthor(tt.nextName, tt.nextEmail)

			assert.Equal(t, tt.wantManagingEditor, p.ManagingEditor)
			assert.Equal(t, tt.wantManagingEditor, p.IAuthor)
			assert.Nil(t, p.IOwner)
		})
	}
}

func TestAddAtomLinkHrefEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)

	// act
	p.AddAtomLink("")

	// assert
	assert.Nil(t, p.AtomLink)
}

func TestAddCategoryEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)

	// act
	p.AddCategory("", nil)

	// assert
	assert.Empty(t, p.ICategories)
	assert.Empty(t, p.Category)
}

func TestAddCategorySubCatEmpty1(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)

	// act
	p.AddCategory("mycat", []string{""})

	// assert
	assert.Len(t, p.ICategories, 1)
	assert.Equal(t, "mycat", p.Category)
	assert.Empty(t, p.ICategories[0].ICategories)
}

func TestAddCategorySubCatEmpty2(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)

	// act
	p.AddCategory("mycat", []string{"xyz", "", "abc"})

	// assert
	assert.Len(t, p.ICategories, 1)
	assert.Equal(t, "mycat", p.Category)
	assert.Len(t, p.ICategories[0].ICategories, 2)
}

func TestAddImageEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)

	// act
	p.AddImage("")

	// assert
	assert.Nil(t, p.Image)
	assert.Nil(t, p.IImage)
}

func TestAddItemEmptyTitleDescription(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 0, added)
	assert.ErrorIs(t, err, podcast.ErrTitleDescriptionRequired)
}

func TestAddItemEmptyEnclosureURL(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.AddEnclosure("", podcast.MP3, 1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 0, added)
	assert.ErrorIs(t, err, podcast.ErrEnclosureURLRequired)
}

func TestAddItemEmptyEnclosureType(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.AddEnclosure("http://example.com/1.mp3", 99, 1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 0, added)
	assert.ErrorIs(t, err, podcast.ErrEnclosureTypeRequired)
}

func TestAddItemZeroValueEnclosureType(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.Enclosure = &podcast.Enclosure{
		URL:    "http://example.com/1.mp3",
		Length: 1,
	}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 0, added)
	assert.ErrorIs(t, err, podcast.ErrEnclosureTypeRequired)
}

func TestAddItemEmptyLink(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc"}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 0, added)
	assert.ErrorIs(t, err, podcast.ErrLinkRequired)
}

func TestSentinelErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		makeErr        func() error
		want           error
		wantValidation bool
		wantTitle      string
	}{
		{
			name: "nil podcast",
			makeErr: func() error {
				var p *podcast.Podcast
				_, err := p.AddItem(podcast.Item{})
				return err
			},
			want: podcast.ErrPodcastRequired,
		},
		{
			name: "nil writer",
			makeErr: func() error {
				p := podcast.New("title", "link", "description", zeroDate, zeroDate)
				return p.Encode(nil)
			},
			want: podcast.ErrWriterRequired,
		},
		{
			name: "missing title and description",
			makeErr: func() error {
				p := podcast.New("title", "link", "description", zeroDate, zeroDate)
				_, err := p.AddItem(podcast.Item{})
				return err
			},
			want:           podcast.ErrTitleDescriptionRequired,
			wantValidation: true,
		},
		{
			name: "missing enclosure url",
			makeErr: func() error {
				p := podcast.New("title", "link", "description", zeroDate, zeroDate)
				i := podcast.Item{Title: "episode", Description: "description"}
				i.AddEnclosure("", podcast.MP3, 1)
				_, err := p.AddItem(i)
				return err
			},
			want:           podcast.ErrEnclosureURLRequired,
			wantValidation: true,
			wantTitle:      "episode",
		},
		{
			name: "missing enclosure type",
			makeErr: func() error {
				p := podcast.New("title", "link", "description", zeroDate, zeroDate)
				i := podcast.Item{Title: "episode", Description: "description"}
				i.AddEnclosure("https://example.com/episode.mp3", podcast.EnclosureUnknown, 1)
				_, err := p.AddItem(i)
				return err
			},
			want:           podcast.ErrEnclosureTypeRequired,
			wantValidation: true,
			wantTitle:      "episode",
		},
		{
			name: "missing link",
			makeErr: func() error {
				p := podcast.New("title", "link", "description", zeroDate, zeroDate)
				_, err := p.AddItem(podcast.Item{Title: "episode", Description: "description"})
				return err
			},
			want:           podcast.ErrLinkRequired,
			wantValidation: true,
			wantTitle:      "episode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.makeErr()
			require.ErrorIs(t, err, tt.want)

			var validationErr *podcast.ItemValidationError
			if got := errors.As(err, &validationErr); got != tt.wantValidation {
				t.Fatalf("errors.As(ItemValidationError) = %v, want %v", got, tt.wantValidation)
			}
			if tt.wantValidation && validationErr.Title != tt.wantTitle {
				t.Fatalf("ItemValidationError.Title = %q, want %q", validationErr.Title, tt.wantTitle)
			}
		})
	}
}

func TestAddItemEnclosureLengthMin(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.AddEnclosure("http://example.com/1.mp3", podcast.MP3, -1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, int64(0), p.Items[0].Enclosure.Length)
}

func TestAddItemEnclosureNoLinkOverride(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.AddEnclosure("http://example.com/1.mp3", podcast.MP3, -1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, i.Enclosure.URL, p.Items[0].Link)
}

func TestAddItemEnclosureLinkPresentNoOverride(t *testing.T) {
	t.Parallel()

	// arrange
	theLink := "http://someotherurl.com/story.html"
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.Link = theLink
	i.AddEnclosure("http://example.com/1.mp3", podcast.MP3, -1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, theLink, p.Items[0].Link)
}

func TestAddItemNoEnclosureGUIDValid(t *testing.T) {
	t.Parallel()

	// arrange
	theLink := "http://someotherurl.com/story.html"
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.Link = theLink

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, theLink, p.Items[0].GUID)
}

func TestAddItemWithEnclosureGUIDSet(t *testing.T) {
	t.Parallel()

	// arrange
	theLink := "http://someotherurl.com/story.html"
	theGUID := "someGUID"
	length := 3
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{
		Title:       "title",
		Description: "desc",
		GUID:        theGUID,
	}
	i.AddEnclosure(theLink, podcast.MP3, int64(length))

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, theGUID, p.Items[0].GUID)
	assert.Equal(t, int64(length), p.Items[0].Enclosure.Length)
}

func TestAddItemAuthor(t *testing.T) {
	t.Parallel()

	// arrange
	theAuthor := podcast.Author{Name: "Jane Doe", Email: "me@janedoe.com"}
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	i := podcast.Item{Title: "title", Description: "desc", Link: "http://a.co/"}
	i.Author = &theAuthor

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, &theAuthor, p.Items[0].Author)
	assert.Equal(t, theAuthor.Email, p.Items[0].IAuthor)
}

func TestAddItemRootManagingEditorSetsAuthorIAuthor(t *testing.T) {
	t.Parallel()

	// arrange
	theAuthor := "me@janedoe.com"
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	p.ManagingEditor = theAuthor
	i := podcast.Item{Title: "title", Description: "desc", Link: "http://a.co/"}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, theAuthor, p.Items[0].Author.Email)
	assert.Equal(t, theAuthor, p.Items[0].IAuthor)
}

func TestAddItemRootIAuthorSetsAuthorIAuthor(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)
	p.IAuthor = "me@janedoe.com"
	i := podcast.Item{Title: "title", Description: "desc", Link: "http://a.co/"}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, "me@janedoe.com", p.Items[0].Author.Email)
	assert.Equal(t, "me@janedoe.com", p.Items[0].IAuthor)
}

func TestAddItemClonesNestedPointers(t *testing.T) {
	t.Parallel()

	enclosure := &podcast.Enclosure{
		URL:    "https://example.com/episode.mp3",
		Type:   podcast.MP3,
		Length: -1,
	}
	author := &podcast.Author{Name: "Jane Doe", Email: "jane@example.com"}
	image := &podcast.IImage{HREF: "https://example.com/episode.png"}
	summary := &podcast.ISummary{Text: "original summary"}
	episodeType := &podcast.IEpisodeType{Text: "full"}
	i := podcast.Item{
		Title:        "episode",
		Description:  "description",
		Enclosure:    enclosure,
		Author:       author,
		IImage:       image,
		ISummary:     summary,
		IEpisodeType: episodeType,
	}
	p := podcast.New("title", "link", "description", zeroDate, zeroDate)

	added, err := p.AddItem(i)

	assert.Equal(t, 1, added)
	require.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.Equal(t, int64(-1), enclosure.Length)
	assert.Empty(t, enclosure.LengthFormatted)
	assert.Equal(t, int64(0), p.Items[0].Enclosure.Length)

	enclosure.URL = "https://example.com/changed.mp3"
	enclosure.Type = podcast.M4A
	enclosure.Length = 99
	author.Email = "changed@example.com"
	image.HREF = "https://example.com/changed.png"
	summary.Text = "changed summary"
	episodeType.Text = "bonus"

	assert.Equal(t, "https://example.com/episode.mp3", p.Items[0].Enclosure.URL)
	assert.Equal(t, podcast.MP3, p.Items[0].Enclosure.Type)
	assert.Equal(t, int64(0), p.Items[0].Enclosure.Length)
	assert.Equal(t, "jane@example.com", p.Items[0].Author.Email)
	assert.Equal(t, "https://example.com/episode.png", p.Items[0].IImage.HREF)
	assert.Equal(t, "original summary", p.Items[0].ISummary.Text)
	assert.Equal(t, "full", p.Items[0].IEpisodeType.Text)
}

func TestAddSubTitleEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "desc", "Link", zeroDate, zeroDate)

	// act
	p.AddSubTitle("")

	// assert
	assert.Empty(t, p.ISubtitle)
}

func TestAddSubTitleTooLong(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "desc", "Link", zeroDate, zeroDate)
	subTitle := ""
	for len(subTitle) < 80 {
		subTitle += "ajd 2 "
	}

	// act
	p.AddSubTitle(subTitle)

	// assert
	assert.Len(t, p.ISubtitle, 64)
}

func TestAddSummaryTooLong(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New(
		"title",
		"desc",
		"Link",
		zeroDate, zeroDate)
	summary := ""
	for len(summary) < 4051 {
		summary += "jax ss 7 "
	}

	// act
	p.AddSummary(summary)

	// assert
	assert.Len(t, p.ISummary.Text, 4000)
}

func TestAddSummaryEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "desc", "Link", zeroDate, zeroDate)

	// act
	p.AddSummary("")

	// assert
	assert.Nil(t, p.ISummary)
}

func TestEncodeSafety(t *testing.T) {
	t.Parallel()

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()

		var p *podcast.Podcast
		var buf bytes.Buffer

		err := p.Encode(&buf)

		assert.ErrorIs(t, err, podcast.ErrPodcastRequired)
	})

	t.Run("nil writer", func(t *testing.T) {
		t.Parallel()

		p := podcast.New("title", "link", "description", zeroDate, zeroDate)

		err := p.Encode(nil)

		assert.ErrorIs(t, err, podcast.ErrWriterRequired)
	})

	t.Run("typed nil writer", func(t *testing.T) {
		t.Parallel()

		p := podcast.New("title", "link", "description", zeroDate, zeroDate)
		var nilBuffer *bytes.Buffer
		var w io.Writer = nilBuffer

		err := p.Encode(w)

		assert.ErrorIs(t, err, podcast.ErrWriterRequired)
	})

	t.Run("zero value uses default encoder", func(t *testing.T) {
		t.Parallel()

		var p podcast.Podcast
		var buf bytes.Buffer

		err := p.Encode(&buf)

		require.NoError(t, err)
		assert.Contains(t, buf.String(), "<rss")
	})
}

func TestStringBytesSafety(t *testing.T) {
	t.Parallel()

	var p *podcast.Podcast

	assert.NotPanics(t, func() {
		_ = p.String()
		_ = p.Bytes()
	})
}

type errWriter struct {
	err error
}

func (w errWriter) Write(_ []byte) (n int, err error) {
	return 0, w.err
}

func TestEncodeWriterError(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "desc", "Link", zeroDate, zeroDate)
	writeErr := errors.New("it was bad")

	// act
	err := p.Encode(&errWriter{err: writeErr})

	// assert
	require.ErrorIs(t, err, writeErr)
	assert.Contains(t, err.Error(), "writing xml header")
}
