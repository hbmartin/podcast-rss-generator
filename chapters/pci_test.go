package chapters_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func expectedPCI() []chapters.Chapter {
	return []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 310, Title: "Main topic", URL: "https://example.com/topic", Image: "https://example.com/topic.png"},
	}
}

func TestExtractPCIChapters(t *testing.T) {
	t.Parallel()
	assert.Equal(t, expectedPCI(), chapters.ExtractPCIChapters(pciJSONDoc(t)))
}

func TestExtractPCIChaptersMissingKey(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractPCIChapters(map[string]any{"nope": []any{}}))
}

func TestExtractPCIChaptersBadStart(t *testing.T) {
	t.Parallel()
	doc := pciJSONDoc(t)
	asMap(t, asSlice(t, doc["chapters"])[0])["startTime"] = "not-a-number"
	assert.Nil(t, chapters.ExtractPCIChapters(doc))
}

func TestGetAndExtractFetches(t *testing.T) {
	t.Parallel()
	client := respondJSON(http.StatusOK, pciJSONDoc(t))
	assert.Equal(t, expectedPCI(), chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json"))
}

func TestGetAndExtractHTTPError(t *testing.T) {
	t.Parallel()
	client := respondJSON(http.StatusInternalServerError, nil)
	assert.Nil(t, chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json"))
}

func TestGetAndExtractRequestError(t *testing.T) {
	t.Parallel()
	client := &fakeClient{do: func(_ *http.Request) (*http.Response, error) {
		return nil, errFakeNetwork
	}}
	assert.Nil(t, chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json"))
}

func TestGetAndExtractBadJSONResponse(t *testing.T) {
	t.Parallel()
	client := respondText(http.StatusOK, "not json at all")
	assert.Nil(t, chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json"))
}

func TestGetAndExtractWritesArchive(t *testing.T) {
	t.Parallel()
	client := respondJSON(http.StatusOK, pciJSONDoc(t))
	archive := filepath.Join(t.TempDir(), "chapters.json")
	got := chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json",
		chapters.WithArchivePath(archive))
	assert.Equal(t, expectedPCI(), got)

	data, err := os.ReadFile(archive) //nolint:gosec // test-controlled path
	require.NoError(t, err)
	var archived map[string]any
	require.NoError(t, json.Unmarshal(data, &archived))
	assert.Equal(t, pciJSONDoc(t), archived)
}

func TestGetAndExtractReturnsChaptersWhenArchiveWriteFails(t *testing.T) {
	t.Parallel()
	client := respondJSON(http.StatusOK, pciJSONDoc(t))
	// A path inside a non-existent directory cannot be written.
	archive := filepath.Join(t.TempDir(), "missing-dir", "chapters.json")
	got := chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json",
		chapters.WithArchivePath(archive))
	assert.Equal(t, expectedPCI(), got)
}

func TestGetAndExtractReadsArchive(t *testing.T) {
	t.Parallel()
	archive := filepath.Join(t.TempDir(), "chapters.json")
	body, err := json.Marshal(pciJSONDoc(t))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(archive, body, 0o600))

	client := &fakeClient{do: func(_ *http.Request) (*http.Response, error) {
		t.Error("should not fetch when archive exists")
		return nil, errNoFetch
	}}
	got := chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json",
		chapters.WithArchivePath(archive))
	assert.Equal(t, expectedPCI(), got)
}

func TestGetAndExtractBadArchiveJSON(t *testing.T) {
	t.Parallel()
	archive := filepath.Join(t.TempDir(), "chapters.json")
	require.NoError(t, os.WriteFile(archive, []byte("{"), 0o600))
	client := respondJSON(http.StatusOK, pciJSONDoc(t))
	assert.Nil(t, chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json",
		chapters.WithArchivePath(archive)))
}

func TestGetAndExtractBadArchiveUTF8(t *testing.T) {
	t.Parallel()
	archive := filepath.Join(t.TempDir(), "chapters.json")
	require.NoError(t, os.WriteFile(archive, []byte{0xff}, 0o600))
	client := respondJSON(http.StatusOK, pciJSONDoc(t))
	assert.Nil(t, chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json",
		chapters.WithArchivePath(archive)))
}

func TestGetAndExtractLogsWhenExtractionFails(t *testing.T) {
	t.Parallel()
	client := respondJSON(http.StatusOK, map[string]any{"version": "1.2.0"})
	assert.Nil(t, chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json"))
}

func TestFindPCIChaptersURL(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "https://example.com/chapters-2.json",
		chapters.FindPCIChaptersURL(writeFeedFile(t), "guid-2"))
}

func TestFindPCIChaptersURLNotDeclared(t *testing.T) {
	t.Parallel()
	feed := writeFeedFile(t)
	assert.Empty(t, chapters.FindPCIChaptersURL(feed, "guid-1"))
	assert.Empty(t, chapters.FindPCIChaptersURL(feed, "no-such-guid"))
}

func TestFindPCIChaptersURLInvalidUTF8(t *testing.T) {
	t.Parallel()
	feed := filepath.Join(t.TempDir(), "feed.xml")
	require.NoError(t, os.WriteFile(feed, []byte{0xff}, 0o600))
	assert.Empty(t, chapters.FindPCIChaptersURL(feed, "guid-1"))
}

func TestFindPCIChaptersURLMissingFile(t *testing.T) {
	t.Parallel()
	assert.Empty(t, chapters.FindPCIChaptersURL(filepath.Join(t.TempDir(), "nope.xml"), "guid-1"))
}

func TestFindPCIChaptersURLUnparseableFeed(t *testing.T) {
	t.Parallel()
	feed := filepath.Join(t.TempDir(), "feed.xml")
	require.NoError(t, os.WriteFile(feed, []byte("not xml at all <<<"), 0o600))
	assert.Empty(t, chapters.FindPCIChaptersURL(feed, "guid-1"))
}
