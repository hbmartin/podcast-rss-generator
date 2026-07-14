package transcript_test

import (
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const podloveXML = `<?xml version="1.0"?>` +
	`<podcast:transcripts xmlns:podcast="http://podlove.org/simple-transcripts">` +
	`<speech><item start="00:00:00.000" end="00:00:01.000">Hello</item></speech>` +
	`</podcast:transcripts>`

const xmlLastBody = `thanks everyone out there who listened to "The Open Source Way". ` +
	`If you enjoyed this episode, please share it and don't miss the next one. ` +
	`We usually publish every last Wednesday of the month, and you'll find us on openSAP ` +
	`and in all those places where you find your other podcasts . ` +
	`Either the mainstream apps that you know, or some of the, themselves, ` +
	`open-source podcast apps. Thanks again and bye bye.`

func TestParsePodloveXML(t *testing.T) {
	t.Parallel()
	tr, err := transcript.ParsePodloveXML(readFixture(t, "Hunting CrossSite Scripting on the Web.xsl"))
	require.NoError(t, err)
	assert.Equal(t, "Welcome to", segBody(t, tr.Segments[0]))
	assert.InDelta(t, 0.82, segStart(t, tr.Segments[0]), 0)
	assert.InDelta(t, 1.66, segEnd(t, tr.Segments[0]), 0)
	assert.Equal(t, xmlLastBody, segBody(t, last(tr.Segments)))
	assert.InDelta(t, 1739.4, segStart(t, last(tr.Segments)), 0)
	assert.InDelta(t, 1766.55, segEnd(t, last(tr.Segments)), 0)
}

func TestParsePodloveXMLDropsEmptySpeechSegments(t *testing.T) {
	t.Parallel()
	xmlString := `<?xml version="1.0"?>
	<podcast:transcripts xmlns:podcast="http://podlove.org/simple-transcripts">
		<speech>
			<item start="00:00:00.000" end="00:00:01.000">Hello</item>
		</speech>
		<speech></speech>
	</podcast:transcripts>
	`
	tr, err := transcript.ParsePodloveXML(xmlString)
	require.NoError(t, err)
	require.Len(t, tr.Segments, 1)
	assert.Equal(t, "Hello", segBody(t, tr.Segments[0]))
	assert.InDelta(t, 0.0, segStart(t, tr.Segments[0]), 0)
	assert.InDelta(t, 1.0, segEnd(t, tr.Segments[0]), 0)
}

func TestParsePodloveXMLNotPodloveRaises(t *testing.T) {
	t.Parallel()
	_, err := transcript.ParsePodloveXML(`<?xml version="1.0"?><root><speech/></root>`)
	var xmlErr *transcript.InvalidXMLError
	require.ErrorAs(t, err, &xmlErr)
}

func TestParsePodloveXMLWithoutTranscriptsTagRaises(t *testing.T) {
	t.Parallel()
	xmlString := `<?xml version="1.0"?>` +
		`<root xmlns:podcast="http://podlove.org/simple-transcripts">` +
		`<other>hi</other></root>`
	_, err := transcript.ParsePodloveXML(xmlString)
	var notFound *transcript.NoTranscriptFoundError
	require.ErrorAs(t, err, &notFound)
}

func TestParsePodloveXMLOnlyEmptySpeechRaises(t *testing.T) {
	t.Parallel()
	xmlString := `<?xml version="1.0"?>` +
		`<podcast:transcripts xmlns:podcast="http://podlove.org/simple-transcripts">` +
		`<speech></speech>` +
		`</podcast:transcripts>`
	_, err := transcript.ParsePodloveXML(xmlString)
	var notFound *transcript.NoTranscriptFoundError
	require.ErrorAs(t, err, &notFound)
}

func TestXMLFileToJSONFileWithMetadata(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "ep.xml")
	writeFile(t, source, podloveXML)
	dest := filepath.Join(t.TempDir(), "ep.json")

	require.NoError(t, transcript.XMLFileToJSONFile(source, dest, map[string]string{"title": "Episode 1"}))

	data := decodeJSONFile(t, dest)
	assert.Equal(t, map[string]any{"title": "Episode 1"}, data["metadata"])
	assert.Equal(t, "Hello", firstSegment(t, data)["body"])
}

func TestXMLFileToJSONFileInvalidWritesNothing(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "bad.xml")
	writeFile(t, source, `<?xml version="1.0"?><root></root>`)
	dest := filepath.Join(t.TempDir(), "out.json")

	err := transcript.XMLFileToJSONFile(source, dest, nil)
	var xmlErr *transcript.InvalidXMLError
	require.ErrorAs(t, err, &xmlErr)
	assert.Contains(t, err.Error(), source)
	assert.NoFileExists(t, dest)
}
