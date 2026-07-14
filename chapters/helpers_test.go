package chapters_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Static errors for fake HTTP clients (err113 forbids inline errors.New).
var (
	errFakeNetwork = errors.New("fake network failure")
	errNoFetch     = errors.New("should not fetch")
)

// mustMarshal JSON-encodes value, panicking on the impossible error.
func mustMarshal(value any) []byte {
	body, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return body
}

// asSlice returns v as a []any, failing the test otherwise.
func asSlice(t *testing.T, v any) []any {
	t.Helper()
	slice, ok := v.([]any)
	require.True(t, ok, "expected []any, got %T", v)
	return slice
}

// asMap returns v as a map[string]any, failing the test otherwise.
func asMap(t *testing.T, v any) map[string]any {
	t.Helper()
	m, ok := v.(map[string]any)
	require.True(t, ok, "expected map[string]any, got %T", v)
	return m
}

// feedXML mirrors the FEED_XML fixture from the Python test suite's conftest.py.
const feedXML = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
     xmlns:psc="http://podlove.org/simple-chapters"
     xmlns:podcast="https://podcastindex.org/namespace/1.0">
  <channel>
    <title>Test Podcast</title>
    <item>
      <title>Episode One</title>
      <guid>guid-1</guid>
      <psc:chapters version="1.2">
        <psc:chapter start="00:00:00" title="Intro"/>
        <psc:chapter start="00:05:10" title="Main topic" href="https://example.com/topic"/>
        <psc:chapter start="01:00:00" title="Outro" image="https://example.com/outro.png"/>
      </psc:chapters>
    </item>
    <item>
      <title>Episode Two</title>
      <guid>guid-2</guid>
      <podcast:chapters url="https://example.com/chapters-2.json" type="application/json+chapters"/>
    </item>
    <item>
      <title>Episode Three</title>
      <guid>guid-3</guid>
    </item>
  </channel>
</rss>
`

// writeFeedFile writes feedXML to a temp file and returns its path.
func writeFeedFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "feed.xml")
	require.NoError(t, os.WriteFile(path, []byte(feedXML), 0o600))
	return path
}

// pciJSONString is the PCI_JSON fixture from conftest.py.
const pciJSONString = `{
  "version": "1.2.0",
  "chapters": [
    {"startTime": 0, "title": "Intro"},
    {"startTime": 310, "title": "Main topic",
     "url": "https://example.com/topic", "img": "https://example.com/topic.png"}
  ]
}`

// pciJSONDoc returns a freshly parsed copy of the PCI_JSON fixture, so tests may
// mutate it in isolation (numbers decode to float64, as real JSON would).
func pciJSONDoc(t *testing.T) map[string]any {
	t.Helper()
	var doc map[string]any
	require.NoError(t, json.Unmarshal([]byte(pciJSONString), &doc))
	return doc
}

// fakeClient is an HTTPClient whose behavior is supplied per test.
type fakeClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (c *fakeClient) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}

// respondJSON returns a fakeClient that always answers with the given status
// and JSON-encoded body.
func respondJSON(status int, value any) *fakeClient {
	body := mustMarshal(value)
	return &fakeClient{do: func(_ *http.Request) (*http.Response, error) {
		return newResponse(status, body), nil
	}}
}

// respondText returns a fakeClient that always answers with the given status
// and raw body.
func respondText(status int, body string) *fakeClient {
	return &fakeClient{do: func(_ *http.Request) (*http.Response, error) {
		return newResponse(status, []byte(body)), nil
	}}
}

func newResponse(status int, body []byte) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     make(http.Header),
	}
}
