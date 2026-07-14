package chapters_test

import (
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	plainText     = "Welcome to the show!\n\n0:00 Intro\n2:30 Interview with Jane\n1:02:03 Wrap-up\n"
	htmlBreaks    = "<p>Chapters:</p><p>(00:00) - Intro<br>(05:10) - The main topic<br>(59:59) - Outro</p>"
	bracketed     = "[0:00] Cold open\n[12:34] Listener questions\n"
	withLink      = "0:00 Intro\n5:00 Sponsor https://example.com/deal\n10:00 Outro\n"
	singleHTMLRun = "<p>0:00 Intro 5:00 Topic two 10:00 The end</p>"
)

func starts(chs []chapters.Chapter) []int {
	out := make([]int, len(chs))
	for i, c := range chs {
		out[i] = c.Start
	}
	return out
}

func TestDescriptionPlainTextLines(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 150, Title: "Interview with Jane"},
		{Start: 3723, Title: "Wrap-up"},
	}, chapters.ExtractDescriptionChapters(plainText, false))
}

func TestDescriptionHTMLWithBreaks(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 310, Title: "The main topic"},
		{Start: 3599, Title: "Outro"},
	}, chapters.ExtractDescriptionChapters(htmlBreaks, false))
}

func TestDescriptionBracketedTimestamps(t *testing.T) {
	t.Parallel()
	chs := chapters.ExtractDescriptionChapters(bracketed, false)
	require.NotNil(t, chs)
	assert.Equal(t, []int{0, 754}, starts(chs))
	assert.Equal(t, "Cold open", chs[0].Title)
}

func TestDescriptionURLExtractedFromTitle(t *testing.T) {
	t.Parallel()
	chs := chapters.ExtractDescriptionChapters(withLink, false)
	require.NotNil(t, chs)
	assert.Equal(t, "https://example.com/deal", chs[1].URL)
	assert.Empty(t, chs[0].URL)
}

func TestDescriptionRetrySplitsSingleRun(t *testing.T) {
	t.Parallel()
	chs := chapters.ExtractDescriptionChapters(singleHTMLRun, false)
	require.NotNil(t, chs)
	assert.Equal(t, []int{0, 300, 600}, starts(chs))
	assert.Equal(t, "Intro", chs[0].Title)
	assert.Equal(t, "Topic two", chs[1].Title)
}

func TestDescriptionInvalidTimestampIsSkipped(t *testing.T) {
	t.Parallel()
	description := "0:00 Intro\n99:99 Bad timestamp\n10:00 The end\n"
	assert.Equal(t, []chapters.Chapter{
		{Start: 0, Title: "Intro"},
		{Start: 600, Title: "The end"},
	}, chapters.ExtractDescriptionChapters(description, false))
}

func TestDescriptionNoChaptersReturnsNil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractDescriptionChapters("Just some show notes.", false))
}

func TestDescriptionSingleTimestampReturnsNil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, chapters.ExtractDescriptionChapters("At 12:30 we talk about X.", false))
}

func TestDescriptionStripHTMLTitles(t *testing.T) {
	t.Parallel()
	description := `0:00 Intro<br>5:00 <a href="https://example.com">A &amp; B</a><br>`
	chs := chapters.ExtractDescriptionChapters(description, true)
	require.NotNil(t, chs)
	assert.Equal(t, "A & B", chs[1].Title)
	assert.Equal(t, "https://example.com", chs[1].URL)
}

func TestDescriptionChapterFieldsAccessibleByName(t *testing.T) {
	t.Parallel()
	chs := chapters.ExtractDescriptionChapters(plainText, false)
	require.NotNil(t, chs)
	assert.Equal(t, 0, chs[0].Start)
	assert.Equal(t, "Intro", chs[0].Title)
	assert.Equal(t, chapters.Chapter{Start: 0, Title: "Intro"}, chs[0])
}
