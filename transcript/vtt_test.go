package transcript_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validVTT = "WEBVTT\n\n00:00.000 --> 00:01.000\nHello from VTT.\n"

func TestParseVTTWithSpeakerTags(t *testing.T) {
	t.Parallel()
	tr := parseVTTFixture(t, "Labs Setting Engineering Goals and Reporting to Stakeholders.txt")
	assert.Equal(t, "1.0.0", tr.Version)
	require.Len(t, tr.Segments, 494)
	assert.Equal(t, "We need to remember that less is more.", segBody(t, tr.Segments[0]))
	assert.InDelta(t, 0.0, segStart(t, tr.Segments[0]), 0)
	assert.InDelta(t, 2.310, segEnd(t, tr.Segments[0]), 0)
	assert.Equal(t, "Yishai Beeri", *segSpeaker(tr.Segments[0]))
	assert.InDelta(t, 2271.601, segStart(t, last(tr.Segments)), 0)
	assert.InDelta(t, 2272.891, segEnd(t, last(tr.Segments)), 0)
	assert.Equal(t, "Talk again soon.", segBody(t, last(tr.Segments)))
	assert.Equal(t, "Dan Lines", *segSpeaker(last(tr.Segments)))
}

func TestParseVTTWithSomeSpeakerAndNote(t *testing.T) {
	t.Parallel()
	tr := parseVTTFixture(t, "Managed services vs. DIY.vtt")
	assert.Equal(t, "1.0.0", tr.Version)
	require.Len(t, tr.Segments, 576)
	assert.InDelta(t, 0.038, segStart(t, tr.Segments[0]), 0)
	assert.InDelta(t, 3.938, segEnd(t, tr.Segments[0]), 0)
	assert.Equal(t, "Hello, and welcome to Postgres FM, a new weekly show", segBody(t, tr.Segments[0]))
	assert.Equal(t, "Michael", *segSpeaker(tr.Segments[0]))
	assert.InDelta(t, 1971.997, segStart(t, last(tr.Segments)), 0)
	assert.InDelta(t, 1972.267, segEnd(t, last(tr.Segments)), 0)
	assert.Equal(t, "Bye.", segBody(t, last(tr.Segments)))
	assert.Equal(t, "Nikolay", *segSpeaker(last(tr.Segments)))
}

func TestParseVTTWithNumberedBlocks(t *testing.T) {
	t.Parallel()
	tr := parseVTTFixture(t, "Zenlytic Is Building You A Better Coworker With AI Agents.txt")
	assert.Equal(t, "1.0.0", tr.Version)
	require.Len(t, tr.Segments, 727)
	assert.InDelta(t, 0.0, segStart(t, tr.Segments[0]), 0)
	assert.InDelta(t, 15.359, segEnd(t, tr.Segments[0]), 0)
	assert.Equal(t,
		"Hello, and welcome to the Data Engineering Podcast, the show about modern data management.",
		segBody(t, tr.Segments[0]))
	assert.Nil(t, segSpeaker(tr.Segments[0]))
	assert.InDelta(t, 3249.44, segStart(t, last(tr.Segments)), 0)
	assert.InDelta(t, 3251.44, segEnd(t, last(tr.Segments)), 0)
	assert.Equal(t, "Podcasts and tell your friends and coworkers.", segBody(t, last(tr.Segments)))
}

func TestParseVTTInvalid(t *testing.T) {
	t.Parallel()
	_, err := transcript.ParseVTT("")
	var vttErr *transcript.InvalidVTTError
	require.ErrorAs(t, err, &vttErr)
}

func TestVTTFileToJSONFileWithMetadata(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "ep.vtt")
	writeFile(t, source, validVTT)
	dest := filepath.Join(t.TempDir(), "ep.json")

	require.NoError(t, transcript.VTTFileToJSONFile(source, dest, map[string]string{"title": "Episode 1"}))

	data := decodeJSONFile(t, dest)
	assert.Equal(t, map[string]any{"title": "Episode 1"}, data["metadata"])
	assert.Equal(t, "Hello from VTT.", firstSegment(t, data)["body"])
}

func TestVTTFileToJSONFileInvalidWritesNothing(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "bad.vtt")
	writeFile(t, source, "this is not a valid vtt file")
	dest := filepath.Join(t.TempDir(), "out.json")

	err := transcript.VTTFileToJSONFile(source, dest, nil)
	var vttErr *transcript.InvalidVTTError
	require.ErrorAs(t, err, &vttErr)
	assert.Contains(t, err.Error(), source)
	assert.NoFileExists(t, dest)
}

func TestVTTFileToJSONFileMissingFile(t *testing.T) {
	t.Parallel()
	dest := filepath.Join(t.TempDir(), "out.json")
	err := transcript.VTTFileToJSONFile(filepath.Join(t.TempDir(), "missing.vtt"), dest, nil)
	require.ErrorIs(t, err, os.ErrNotExist)
	assert.NoFileExists(t, dest)
}

func parseVTTFixture(t *testing.T, name string) *transcript.Transcript {
	t.Helper()
	tr, err := transcript.ParseVTT(readFixture(t, name))
	require.NoError(t, err)
	return tr
}
