package podcast_test

import (
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newValidItem() podcast.Item {
	return podcast.Item{
		Title:       "Episode 1",
		Description: "Description for Episode 1",
		Link:        "http://example.com/1",
	}
}

func TestAddContentEncodedItem(t *testing.T) {
	t.Parallel()

	i := newValidItem()

	// empty input is a no-op.
	i.AddContentEncoded("")
	assert.Nil(t, i.ContentEncoded)

	// non-empty input populates the field.
	i.AddContentEncoded(`<p>Show notes with a <a href="https://example.com">link</a></p>`)
	require.NotNil(t, i.ContentEncoded)
	assert.Equal(
		t,
		`<p>Show notes with a <a href="https://example.com">link</a></p>`,
		i.ContentEncoded.Text,
	)
}

func TestContentEncodedRendersCDATAAndNamespace(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)

	// Without any content:encoded item, the namespace is not declared.
	assert.NotContains(t, p.String(), "xmlns:content")

	item := newValidItem()
	item.AddContentEncoded(`<p>Rich <b>HTML</b> body & more</p>`)
	_, err := p.AddItem(item)
	require.NoError(t, err)

	got := p.String()
	// Namespace is declared on <rss> once an item uses it.
	assert.Contains(t, got, `xmlns:content="http://purl.org/rss/1.0/modules/content/"`)
	// The body is wrapped in CDATA so raw HTML (and &) is preserved verbatim.
	assert.Contains(
		t,
		got,
		"<content:encoded><![CDATA[<p>Rich <b>HTML</b> body & more</p>]]></content:encoded>",
	)
}

func TestContentEncodedSnapshotAfterAddItem(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)

	item := newValidItem()
	item.AddContentEncoded("<p>original</p>")
	_, err := p.AddItem(item)
	require.NoError(t, err)

	// Mutating the caller's item after AddItem must not change the stored feed.
	item.ContentEncoded.Text = "<p>mutated</p>"

	got := p.String()
	assert.Contains(t, got, "<![CDATA[<p>original</p>]]>")
	assert.NotContains(t, got, "mutated")
}

func TestSetExplicitPodcast(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)

	require.Nil(t, p.IExplicit)

	p.SetExplicit(true)
	assert.Contains(t, p.String(), "<itunes:explicit>true</itunes:explicit>")

	p.SetExplicit(false)
	assert.Contains(t, p.String(), "<itunes:explicit>false</itunes:explicit>")
}

func TestSetCompleteAndBlockPodcast(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)

	// false leaves the tags omitted.
	p.SetComplete(false)
	p.SetBlock(false)
	got := p.String()
	assert.NotContains(t, got, "itunes:complete")
	assert.NotContains(t, got, "itunes:block")

	// true renders the literal "Yes".
	p.SetComplete(true)
	p.SetBlock(true)
	got = p.String()
	assert.Contains(t, got, "<itunes:complete>Yes</itunes:complete>")
	assert.Contains(t, got, "<itunes:block>Yes</itunes:block>")

	// flipping back to false clears the tags again.
	p.SetComplete(false)
	p.SetBlock(false)
	got = p.String()
	assert.NotContains(t, got, "itunes:complete")
	assert.NotContains(t, got, "itunes:block")
}

func TestSetExplicitAndBlockItem(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)
	item := newValidItem()
	item.SetExplicit(true)
	item.SetBlock(true)

	_, err := p.AddItem(item)
	require.NoError(t, err)

	got := p.String()
	assert.Contains(t, got, "<itunes:explicit>true</itunes:explicit>")
	assert.Contains(t, got, "<itunes:block>Yes</itunes:block>")
}

func TestAddItemSnapshotsExplicitAndBlock(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)
	item := newValidItem()
	item.SetExplicit(true)
	item.SetBlock(true)

	_, err := p.AddItem(item)
	require.NoError(t, err)

	// Mutating the caller's Item after AddItem must not affect the stored copy.
	item.SetExplicit(false)
	item.SetBlock(false)

	got := p.String()
	assert.Contains(t, got, "<itunes:explicit>true</itunes:explicit>")
	assert.Contains(t, got, "<itunes:block>Yes</itunes:block>")
}

func TestAddItems(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)

	first := newValidItem()
	second := newValidItem()
	second.Title = "Episode 2"
	second.Link = "http://example.com/2"

	count, err := p.AddItems(first, second)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, p.Items, 2)
}

func TestAddItemsStopsAtFirstInvalid(t *testing.T) {
	t.Parallel()

	p := podcast.New("t", "l", "d", zeroDate, zeroDate)

	good := newValidItem()
	bad := podcast.Item{Title: "no description"} // missing Description

	count, err := p.AddItems(good, bad)
	require.Error(t, err)
	require.ErrorIs(t, err, podcast.ErrTitleDescriptionRequired)
	// The valid item before the failure is retained.
	assert.Equal(t, 1, count)
	assert.Len(t, p.Items, 1)
}

func TestAddItemsNilReceiver(t *testing.T) {
	t.Parallel()

	var p *podcast.Podcast
	_, err := p.AddItems(newValidItem())
	require.Error(t, err)
	assert.ErrorIs(t, err, podcast.ErrPodcastRequired)
}

func TestValidateValidFeed(t *testing.T) {
	t.Parallel()

	p := podcast.New("Title", "http://example.com", "Description", zeroDate, zeroDate)
	_, err := p.AddItem(newValidItem())
	require.NoError(t, err)

	assert.NoError(t, p.Validate())
}

func TestValidateReportsAllChannelErrors(t *testing.T) {
	t.Parallel()

	// Zero-value podcast: no title, description, or link.
	var p podcast.Podcast

	err := p.Validate()
	require.Error(t, err)
	// Validate joins every problem, so assert (not require) on each so all
	// three are checked even if one is missing.
	//nolint:testifylint // intentionally reporting all joined errors, not failing fast
	assert.ErrorIs(t, err, podcast.ErrChannelTitleRequired)
	//nolint:testifylint // intentionally reporting all joined errors, not failing fast
	assert.ErrorIs(t, err, podcast.ErrChannelDescriptionRequired)
	assert.ErrorIs(t, err, podcast.ErrChannelLinkRequired)
}

func TestValidateReportsItemErrors(t *testing.T) {
	t.Parallel()

	p := podcast.New("Title", "http://example.com", "Description", zeroDate, zeroDate)
	// Bypass AddItem's validation by attaching an invalid item directly.
	p.Items = append(p.Items, &podcast.Item{Title: "missing description"})

	err := p.Validate()
	require.Error(t, err)
	require.ErrorIs(t, err, podcast.ErrTitleDescriptionRequired)
	assert.Contains(t, err.Error(), "missing description")
}

func TestValidateNilReceiver(t *testing.T) {
	t.Parallel()

	var p *podcast.Podcast
	err := p.Validate()
	require.Error(t, err)
	assert.ErrorIs(t, err, podcast.ErrPodcastRequired)
}

func TestValidateSkipsNilItems(t *testing.T) {
	t.Parallel()

	p := podcast.New("Title", "http://example.com", "Description", zeroDate, zeroDate)
	p.Items = append(p.Items, nil)

	assert.NoError(t, p.Validate())
}

func TestExplicitFalseIsSerialized(t *testing.T) {
	t.Parallel()

	// Regression: a false explicit flag must render "false", not be omitted,
	// because Apple uses it to display the "Clean" badge.
	p := podcast.New("t", "l", "d", zeroDate, zeroDate)
	p.SetExplicit(false)

	assert.Contains(t, p.String(), "<itunes:explicit>false</itunes:explicit>")
}
