package transcript_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const html300Seg0Body = "Welcome to Go Time. This is a very special episode. This is episode number 300. " +
	"So today we're doing something a little different from our usual content. " +
	"We're having a full panel episode with our co-hosts, with all of our hosts... " +
	"And we're talking about the past of Go Time, the current present, and our plans for the future. " +
	"So joining me live is Jon Calhoun. How are you doing, Jon?"

func TestParseHTMLNoTime(t *testing.T) {
	t.Parallel()
	tr := parseHTMLFixture(t, "300 multiple choices.html")
	assert.Equal(t, "1.0.0", tr.Version)
	require.Len(t, tr.Segments, 307)
	assert.Equal(t, html300Seg0Body, segBody(t, tr.Segments[0]))
	assert.Equal(t, "Kris Brandow", *segSpeaker(tr.Segments[0]))
	assert.Equal(t, "Kris Brandow", *segSpeaker(last(tr.Segments)))
}

func TestParseHTMLNoBody(t *testing.T) {
	t.Parallel()
	tr := parseHTMLFixture(t, "Talking AI at OpenShift Commons Gathering in Raleigh.html")
	assert.Equal(t, "1.0.0", tr.Version)
	require.Len(t, tr.Segments, 131)
	assert.Equal(t,
		"Now for something completely different. Frank is in Raleigh speaking at",
		segBody(t, tr.Segments[0]))
	assert.Equal(t, "Speaker", *segSpeaker(tr.Segments[0]))
	assert.InDelta(t, 0.0, segStart(t, tr.Segments[0]), 0)
	assert.Equal(t, "Thanks. And you", segBody(t, last(tr.Segments)))
	assert.Equal(t, "Speaker", *segSpeaker(last(tr.Segments)))
	assert.InDelta(t, 511.0, segStart(t, last(tr.Segments)), 0)
}

func TestParseHTMLNoCiteOrTime(t *testing.T) {
	t.Parallel()
	_, err := transcript.ParseHTML("<html><body><p>Just a paragraph</p></body></html>")
	var htmlErr *transcript.InvalidHTMLError
	require.ErrorAs(t, err, &htmlErr)
}

func TestParseHTMLWithTimestampInPTag(t *testing.T) {
	t.Parallel()
	tr := parseHTMLFixture(t, "78 Exploring MCMC Sampler Algorithms, with Matt D. Hoffman.html")
	require.Len(t, tr.Segments, 104)
	assert.Equal(t,
		"Okay, I mean now, can you hear me well yes. Okay, can I hear you don't need to",
		segBody(t, tr.Segments[0]))
	assert.InDelta(t, 119.0, segStart(t, tr.Segments[0]), 0)
	assert.Equal(t, "Thank you. You too. Bye.", segBody(t, last(tr.Segments)))
	assert.InDelta(t, 5254.0, segStart(t, last(tr.Segments)), 0)
}

func TestParseHTMLWithTimestampLookingInBody(t *testing.T) {
	t.Parallel()
	tr := parseHTMLFixture(t, "Tales from Manufacturing Shipping Rack 1.html")
	assert.Equal(t, "Speaker 1", *segSpeaker(tr.Segments[0]))
	assert.InDelta(t, 0.0, segStart(t, tr.Segments[0]), 0)
	assert.Equal(t, "Speaker 2", *segSpeaker(last(tr.Segments)))
	assert.InDelta(t, 5051.0, segStart(t, last(tr.Segments)), 0)
}

func TestParseHTMLEmpty(t *testing.T) {
	t.Parallel()
	_, err := transcript.ParseHTML("")
	var htmlErr *transcript.InvalidHTMLError
	require.ErrorAs(t, err, &htmlErr)
}

func TestParseHTMLNoUsableTranscriptContent(t *testing.T) {
	t.Parallel()
	// "<time>" as a substring passes the format gate, but the parsed document
	// holds no cite/time/p content, so no segment is populated.
	_, err := transcript.ParseHTML("<!-- <time> --><div>hello</div>")
	var notFound *transcript.NoTranscriptFoundError
	require.ErrorAs(t, err, &notFound)
}

func TestHTMLFileToJSONFileWithMetadata(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "ep.html")
	writeFile(t, source, "<html><body><cite>Ann:</cite><p>Hello there</p></body></html>")
	dest := filepath.Join(t.TempDir(), "ep.json")

	require.NoError(t, transcript.HTMLFileToJSONFile(source, dest, map[string]string{"title": "Episode 1"}))

	data := decodeJSONFile(t, dest)
	assert.Equal(t, map[string]any{"title": "Episode 1"}, data["metadata"])
	segment := firstSegment(t, data)
	assert.Equal(t, "Ann", segment["speaker"])
	assert.Equal(t, "Hello there", segment["body"])
}

func TestHTMLFileToJSONFileInvalidWritesNothing(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "bad.html")
	writeFile(t, source, "<html><body><p>no cite or time here</p></body></html>")
	dest := filepath.Join(t.TempDir(), "out.json")

	err := transcript.HTMLFileToJSONFile(source, dest, nil)
	var htmlErr *transcript.InvalidHTMLError
	require.ErrorAs(t, err, &htmlErr)
	assert.Contains(t, err.Error(), source)
	assert.NoFileExists(t, dest)
}

func TestHTMLFileToJSONFileMissingFile(t *testing.T) {
	t.Parallel()
	dest := filepath.Join(t.TempDir(), "out.json")
	err := transcript.HTMLFileToJSONFile(filepath.Join(t.TempDir(), "missing.html"), dest, nil)
	require.ErrorIs(t, err, os.ErrNotExist)
	assert.NoFileExists(t, dest)
}

func parseHTMLFixture(t *testing.T, name string) *transcript.Transcript {
	t.Helper()
	tr, err := transcript.ParseHTML(readFixture(t, name))
	require.NoError(t, err)
	return tr
}
