package chapters_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func expectedPSC() []chapters.Chapter {
	return []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 310, Title: "Main topic", URL: "https://example.com/topic"},
		{Start: 3600, Title: "Outro", Image: "https://example.com/outro.png"},
	}
}

func TestExtractByGUID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, expectedPSC(), chapters.ExtractPSCChaptersFromFile(writeFeedFile(t), "guid-1"))
}

func TestExtractMissingGUID(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractPSCChaptersFromFile(writeFeedFile(t), "no-such-guid"))
}

func TestExtractEpisodeWithoutChapters(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractPSCChaptersFromFile(writeFeedFile(t), "guid-3"))
}

func TestExtractMissingFile(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractPSCChaptersFromFile(filepath.Join(t.TempDir(), "nope.xml"), "guid-1"))
}

func TestExtractUnparseableFeed(t *testing.T) {
	t.Parallel()
	bad := filepath.Join(t.TempDir(), "bad.xml")
	require.NoError(t, os.WriteFile(bad, []byte("not xml at all <<<"), 0o600))
	assert.Nil(t, chapters.ExtractPSCChaptersFromFile(bad, "guid-1"))
}

func TestExtractInvalidUTF8FeedReturnsNil(t *testing.T) {
	t.Parallel()
	bad := filepath.Join(t.TempDir(), "bad.xml")
	require.NoError(t, os.WriteFile(bad, []byte{0xff}, 0o600))
	assert.Nil(t, chapters.ExtractPSCChaptersFromFile(bad, "guid-1"))
	assert.Nil(t, chapters.ExtractAllPSCChaptersFromFile(bad))
}

func TestExtractFeedWithoutChannel(t *testing.T) {
	t.Parallel()
	bad := filepath.Join(t.TempDir(), "nochannel.xml")
	require.NoError(t, os.WriteFile(bad, []byte("<rss></rss>"), 0o600))
	assert.Nil(t, chapters.ExtractPSCChaptersFromFile(bad, "guid-1"))
}

func TestExtractAllEpisodes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, map[string][]chapters.Chapter{"guid-1": expectedPSC()},
		chapters.ExtractAllPSCChaptersFromFile(writeFeedFile(t)))
}

func TestExtractAllMissingFile(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractAllPSCChaptersFromFile(filepath.Join(t.TempDir(), "nope.xml")))
}

func TestExtractAllUnparseableFeed(t *testing.T) {
	t.Parallel()
	bad := filepath.Join(t.TempDir(), "bad.xml")
	require.NoError(t, os.WriteFile(bad, []byte("not xml at all <<<"), 0o600))
	assert.Nil(t, chapters.ExtractAllPSCChaptersFromFile(bad))
}

func TestExtractAllSkipsItemWithoutGUID(t *testing.T) {
	t.Parallel()
	feed := filepath.Join(t.TempDir(), "feed.xml")
	content := `<rss xmlns:psc="http://podlove.org/simple-chapters"><channel>` +
		`<item><title>No guid</title>` +
		`<psc:chapters><psc:chapter start="0:00" title="A"/>` +
		`<psc:chapter start="1:00" title="B"/></psc:chapters></item>` +
		`<item><title>Has guid</title><guid>guid-1</guid>` +
		`<psc:chapters><psc:chapter start="0:00" title="Intro"/>` +
		`<psc:chapter start="5:10" title="Main topic" ` +
		`href="https://example.com/topic"/>` +
		`<psc:chapter start="01:00:00" title="Outro" ` +
		`image="https://example.com/outro.png"/></psc:chapters></item>` +
		`</channel></rss>`
	require.NoError(t, os.WriteFile(feed, []byte(content), 0o600))
	assert.Equal(t, map[string][]chapters.Chapter{"guid-1": expectedPSC()},
		chapters.ExtractAllPSCChaptersFromFile(feed))
}

func TestExtractFromElementWithBadStart(t *testing.T) {
	t.Parallel()
	element, err := chapters.ParseXML([]byte(`<chapters><chapter start="bogus" title="X"/></chapters>`))
	require.NoError(t, err)
	assert.Nil(t, chapters.ExtractPSCChapters(element))
}

func TestExtractFromElementMissingTitle(t *testing.T) {
	t.Parallel()
	element, err := chapters.ParseXML([]byte(`<chapters><chapter start="0:10"/></chapters>`))
	require.NoError(t, err)
	assert.Nil(t, chapters.ExtractPSCChapters(element))
}

func TestExtractFromURL(t *testing.T) {
	t.Parallel()
	var capturedURL string
	var hadDeadline bool
	client := &fakeClient{do: func(req *http.Request) (*http.Response, error) {
		capturedURL = req.URL.String()
		_, hadDeadline = req.Context().Deadline()
		return newResponse(http.StatusOK, []byte(feedXML)), nil
	}}
	got := chapters.ExtractPSCChaptersFromURL(context.Background(), client, "https://example.com/feed.xml", "guid-1")
	assert.Equal(t, expectedPSC(), got)
	assert.Equal(t, "https://example.com/feed.xml", capturedURL)
	assert.True(t, hadDeadline, "request should carry a timeout deadline")
}

func TestExtractFromURLUnparseableFeed(t *testing.T) {
	t.Parallel()
	client := respondText(http.StatusOK, "not xml at all <<<")
	assert.Nil(t, chapters.ExtractPSCChaptersFromURL(
		context.Background(), client, "https://example.com/feed.xml", "guid-1"))
}

func TestExtractFromURLHTTPError(t *testing.T) {
	t.Parallel()
	client := respondText(http.StatusNotFound, "")
	assert.Nil(t, chapters.ExtractPSCChaptersFromURL(
		context.Background(), client, "https://example.com/feed.xml", "guid-1"))
}

func TestExtractFromURLRequestError(t *testing.T) {
	t.Parallel()
	client := &fakeClient{do: func(_ *http.Request) (*http.Response, error) {
		return nil, errFakeNetwork
	}}
	assert.Nil(t, chapters.ExtractPSCChaptersFromURL(
		context.Background(), client, "https://example.com/feed.xml", "guid-1"))
}
