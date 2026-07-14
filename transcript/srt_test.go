package transcript_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validSRT = "1\n00:00:00,000 --> 00:00:01,000\nMichael: Hello world.\n\n"

func TestParseSRT(t *testing.T) {
	t.Parallel()
	tr := parseSRTFixture(t, "Bloat.srt")
	assert.Equal(t, "1.0.0", tr.Version)
	require.Len(t, tr.Segments, 724)
	assert.Equal(t, "Michael", *segSpeaker(tr.Segments[0]))
	assert.Equal(t, "Hello, and welcome to PostgresFM, a weekly show about", segBody(t, tr.Segments[0]))
	assert.InDelta(t, 0.060, segStart(t, tr.Segments[0]), 0)
	assert.InDelta(t, 2.680, segEnd(t, tr.Segments[0]), 0)
	assert.Equal(t, "Take care, Chelsea.", segBody(t, last(tr.Segments)))
	assert.InDelta(t, 2173.940, segStart(t, last(tr.Segments)), 0)
	assert.InDelta(t, 2174.920, segEnd(t, last(tr.Segments)), 0)
}

func TestParseSRTWithExtraWhitespace(t *testing.T) {
	t.Parallel()
	tr := parseSRTFixture(t, "I Quit My Job.srt")
	assert.Equal(t, "Welcome to another episode of the map escaping podcast.", segBody(t, tr.Segments[0]))
	assert.Equal(t, "Bye.", segBody(t, last(tr.Segments)))
}

func TestParseSRTWithMissingWhitespace(t *testing.T) {
	t.Parallel()
	tr := parseSRTFixture(t, "AI Autocomplete for QGIS.srt")
	assert.Equal(t, "Hi,Brendan.Welcometothepodcast.", segBody(t, tr.Segments[0]))
	assert.Equal(t, "Thanks.", segBody(t, last(tr.Segments)))
}

func TestParseSRTWithPeriodTimestamps(t *testing.T) {
	t.Parallel()
	tr := parseSRTFixture(t, "Episode 17.srt")
	assert.InDelta(t, 5.0, segEnd(t, tr.Segments[0]), 0)
	assert.Equal(t,
		"I think in terms of news, what would be great to kick off with is Swift Server.",
		segBody(t, tr.Segments[0]))
	assert.InDelta(t, 2882.360, segEnd(t, last(tr.Segments)), 0)
	assert.Equal(t, "Yeah.", segBody(t, last(tr.Segments)))
}

func TestParseSRTWithNewlinesInBody(t *testing.T) {
	t.Parallel()
	tr := parseSRTFixture(t, "Yak Shaving with Tim Mitra.srt")
	assert.Equal(t, "Leo Dion (host)", *segSpeaker(tr.Segments[0]))
	assert.Equal(t,
		"Before we begin today's episode, I wanted to let you know Bright Digit needs your help.",
		segBody(t, tr.Segments[0]))
	assert.Equal(t, "?   Yeah.", segBody(t, tr.Segments[881]))
	assert.Equal(t, "Bye everyone.", segBody(t, last(tr.Segments)))
}

func TestParseSRTInvalid(t *testing.T) {
	t.Parallel()
	_, err := transcript.ParseSRT("whatever")
	var srtErr *transcript.InvalidSRTError
	require.ErrorAs(t, err, &srtErr)
}

func TestSRTFileToJSONFileWithMetadata(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "ep.srt")
	writeFile(t, source, validSRT)
	dest := filepath.Join(t.TempDir(), "ep.json")

	require.NoError(t, transcript.SRTFileToJSONFile(source, dest, map[string]string{"title": "Episode 1"}))

	data := decodeJSONFile(t, dest)
	assert.Equal(t, map[string]any{"title": "Episode 1"}, data["metadata"])
	segment := firstSegment(t, data)
	assert.Equal(t, "Michael", segment["speaker"])
	assert.Equal(t, "Hello world.", segment["body"])
}

func TestSRTFileToJSONFileInvalidWritesNothing(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "bad.srt")
	writeFile(t, source, "this is not valid srt")
	dest := filepath.Join(t.TempDir(), "out.json")

	err := transcript.SRTFileToJSONFile(source, dest, nil)
	var srtErr *transcript.InvalidSRTError
	require.ErrorAs(t, err, &srtErr)
	assert.Contains(t, err.Error(), source)
	assert.NoFileExists(t, dest)
}

func parseSRTFixture(t *testing.T, name string) *transcript.Transcript {
	t.Helper()
	tr, err := transcript.ParseSRT(readFixture(t, name))
	require.NoError(t, err)
	return tr
}

func decodeJSONFile(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path) //nolint:gosec // test-controlled path
	require.NoError(t, err)
	var out map[string]any
	require.NoError(t, json.Unmarshal(data, &out))
	return out
}
