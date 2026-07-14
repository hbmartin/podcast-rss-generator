package transcript_test

import (
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/require"
)

// testdataPath returns the path to a named fixture under testdata.
func testdataPath(name string) string {
	return filepath.Join("testdata", name)
}

// readFixture reads a fixture the way the converters do: as UTF-8 text with
// universal-newline translation, mirroring Python's Path.read_text.
func readFixture(t *testing.T, name string) string {
	t.Helper()
	content, err := transcript.ReadTextRobust(testdataPath(name))
	require.NoError(t, err)
	return content
}

func segBody(t *testing.T, seg transcript.Segment) string {
	t.Helper()
	require.NotNil(t, seg.Body)
	return *seg.Body
}

func segSpeaker(seg transcript.Segment) *string {
	return seg.Speaker
}

func segStart(t *testing.T, seg transcript.Segment) float64 {
	t.Helper()
	require.NotNil(t, seg.StartTime)
	return *seg.StartTime
}

func segEnd(t *testing.T, seg transcript.Segment) float64 {
	t.Helper()
	require.NotNil(t, seg.EndTime)
	return *seg.EndTime
}

// last returns the final segment of a transcript.
func last(segments []transcript.Segment) transcript.Segment {
	return segments[len(segments)-1]
}

// segmentsOf returns the decoded "segments" array of a transcript JSON object.
func segmentsOf(t *testing.T, data map[string]any) []any {
	t.Helper()
	segments, ok := data["segments"].([]any)
	require.True(t, ok, "segments should be an array")
	return segments
}

// firstSegment returns the first decoded segment object of a transcript JSON.
func firstSegment(t *testing.T, data map[string]any) map[string]any {
	t.Helper()
	segment, ok := segmentsOf(t, data)[0].(map[string]any)
	require.True(t, ok, "segment should be an object")
	return segment
}
