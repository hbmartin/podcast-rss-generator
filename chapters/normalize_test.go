package chapters_test

import (
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
)

func titles(chs []chapters.Chapter) []string {
	out := make([]string, len(chs))
	for i, c := range chs {
		out[i] = c.Title
	}
	return out
}

func TestStripHTML(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "A & B", chapters.StripHTML(`<a href="https://x.com">A &amp; B</a>`))
	assert.Equal(t, "plain title", chapters.StripHTML("plain title"))
	assert.Equal(t, "a b", chapters.StripHTML("a<br>b"))
}

func TestNormalizeSortsByStart(t *testing.T) {
	t.Parallel()
	chs := []chapters.Chapter{{Start: 60, Title: "b"}, {Start: 0, Title: "a"}, {Start: 30, Title: "c"}}
	assert.Equal(t, []string{"a", "c", "b"}, titles(chapters.NormalizeChapters(chs)))
}

func TestNormalizeSortDisabled(t *testing.T) {
	t.Parallel()
	chs := []chapters.Chapter{{Start: 60, Title: "b"}, {Start: 0, Title: "a"}}
	assert.Equal(t, []string{"b", "a"}, titles(chapters.NormalizeChapters(chs, chapters.WithoutSort())))
}

func TestNormalizeDedupesRepeatedStarts(t *testing.T) {
	t.Parallel()
	chs := []chapters.Chapter{{Start: 0, Title: "a"}, {Start: 0, Title: "dup"}, {Start: 10, Title: "b"}}
	assert.Equal(t, []string{"a", "b"}, titles(chapters.NormalizeChapters(chs)))
}

func TestNormalizeDedupeDisabled(t *testing.T) {
	t.Parallel()
	chs := []chapters.Chapter{{Start: 0, Title: "a"}, {Start: 0, Title: "dup"}}
	assert.Len(t, chapters.NormalizeChapters(chs, chapters.WithoutDedupe()), 2)
}

func TestNormalizeClampsToMaxStart(t *testing.T) {
	t.Parallel()
	chs := []chapters.Chapter{{Start: 0, Title: "a"}, {Start: 5000, Title: "past-the-end"}}
	assert.Equal(t, []string{"a"}, titles(chapters.NormalizeChapters(chs, chapters.WithMaxStart(3600))))
}

func TestNormalizeDropsNegativeStarts(t *testing.T) {
	t.Parallel()
	chs := []chapters.Chapter{{Start: -5, Title: "bad"}, {Start: 0, Title: "a"}}
	assert.Equal(t, []string{"a"}, titles(chapters.NormalizeChapters(chs)))
}

func TestNormalizeStripTitles(t *testing.T) {
	t.Parallel()
	chs := []chapters.Chapter{{Start: 0, Title: "<b>Intro</b>"}}
	assert.Equal(t, "Intro", chapters.NormalizeChapters(chs, chapters.WithStripTitles())[0].Title)
}

func TestNormalizeResultOrdering(t *testing.T) {
	t.Parallel()
	result := chapters.NormalizeChapters([]chapters.Chapter{{Start: 10, Title: "b"}, {Start: 0, Title: "a"}})
	assert.Equal(t, []chapters.Chapter{{Start: 0, Title: "a"}, {Start: 10, Title: "b"}}, result)
}
