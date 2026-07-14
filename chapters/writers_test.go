package chapters_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writerChapters() []chapters.Chapter {
	return []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 310, Title: "Main topic", URL: "https://example.com/topic"},
		{Start: 3722, Title: "Outro", Image: "https://example.com/outro.png"},
	}
}

func TestToPCIDict(t *testing.T) {
	t.Parallel()
	doc := chapters.ToPCIDict(writerChapters())
	assert.Equal(t, "1.2.0", doc.Version)
	assert.Equal(t, []chapters.PCIChapter{
		{StartTime: 0, Title: "Intro"},
		{StartTime: 310, Title: "Main topic", URL: "https://example.com/topic"},
		{StartTime: 3722, Title: "Outro", Img: "https://example.com/outro.png"},
	}, doc.Chapters)
}

func TestToPCIJSONRoundtrip(t *testing.T) {
	t.Parallel()
	rendered, err := chapters.ToPCIJSON(writerChapters(), 2)
	require.NoError(t, err)

	var doc map[string]any
	require.NoError(t, json.Unmarshal([]byte(rendered), &doc))
	assert.Equal(t, []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 310, Title: "Main topic", URL: "https://example.com/topic"},
		{Start: 3722, Title: "Outro", Image: "https://example.com/outro.png"},
	}, chapters.ExtractPCIChapters(doc))
}

func TestToPSCXMLRoundtrip(t *testing.T) {
	t.Parallel()
	rendered, err := chapters.ToPSCXML(writerChapters())
	require.NoError(t, err)
	assert.Contains(t, rendered, `xmlns:psc="http://podlove.org/simple-chapters"`)

	element, err := chapters.ParseXML([]byte(rendered))
	require.NoError(t, err)
	assert.Equal(t, []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 310, Title: "Main topic", URL: "https://example.com/topic"},
		{Start: 3722, Title: "Outro", Image: "https://example.com/outro.png"},
	}, chapters.ExtractPSCChapters(element))
}

func TestToPSCElementUsesNamespaceExpandedTags(t *testing.T) {
	t.Parallel()
	element, err := chapters.ToPSCElement(writerChapters())
	require.NoError(t, err)
	assert.Equal(t, "{http://podlove.org/simple-chapters}chapters", element.Tag)
	assert.Equal(t, "{http://podlove.org/simple-chapters}chapter", element.Children[0].Tag)
}

func TestToPSCXMLEscapesTitles(t *testing.T) {
	t.Parallel()
	rendered, err := chapters.ToPSCXML([]chapters.Chapter{{Start: 0, Title: `Q&A "special" <session>`}})
	require.NoError(t, err)

	element, err := chapters.ParseXML([]byte(rendered))
	require.NoError(t, err)
	extracted := chapters.ExtractPSCChapters(element)
	require.NotNil(t, extracted)
	assert.Equal(t, `Q&A "special" <session>`, extracted[0].Title)
}

func TestToDescriptionOutput(t *testing.T) {
	t.Parallel()
	rendered, err := chapters.ToDescription(writerChapters())
	require.NoError(t, err)
	assert.Equal(t, []string{
		"0:00 Intro",
		"5:10 Main topic https://example.com/topic",
		"1:02:02 Outro",
	}, strings.Split(rendered, "\n"))
}

func TestToDescriptionRoundtrip(t *testing.T) {
	t.Parallel()
	rendered, err := chapters.ToDescription(writerChapters())
	require.NoError(t, err)

	extracted := chapters.ExtractDescriptionChapters(rendered, false)
	require.NotNil(t, extracted)
	starts := make([]int, len(extracted))
	for i, c := range extracted {
		starts[i] = c.Start
	}
	assert.Equal(t, []int{0, 310, 3722}, starts)
	assert.Equal(t, []string{
		"Intro",
		"Main topic https://example.com/topic",
		"Outro",
	}, titles(extracted))
	assert.Equal(t, "https://example.com/topic", extracted[1].URL)
}
