package chapters

import (
	"html"
	"sort"
	"strings"
)

// StripHTML removes HTML tags and unescapes entities, collapsing runs of
// whitespace to single spaces and trimming the result.
func StripHTML(text string) string {
	replaced := htmlTagRe.ReplaceAllString(text, " ")
	unescaped := html.UnescapeString(replaced)
	return strings.Join(strings.Fields(unescaped), " ")
}

// normalizeConfig holds the resolved options for NormalizeChapters.
type normalizeConfig struct {
	sort        bool
	dedupe      bool
	maxStart    *int
	stripTitles bool
}

// NormalizeOption configures NormalizeChapters. By default chapters are sorted
// by start time and duplicate start times are dropped.
type NormalizeOption func(*normalizeConfig)

// WithoutSort keeps chapters in their original order instead of sorting by
// start time.
func WithoutSort() NormalizeOption {
	return func(c *normalizeConfig) { c.sort = false }
}

// WithoutDedupe keeps chapters that repeat an earlier start time.
func WithoutDedupe() NormalizeOption {
	return func(c *normalizeConfig) { c.dedupe = false }
}

// WithMaxStart drops chapters starting after maxStart seconds (for example the
// episode duration).
func WithMaxStart(maxStart int) NormalizeOption {
	return func(c *normalizeConfig) { c.maxStart = &maxStart }
}

// WithStripTitles removes HTML markup from chapter titles.
func WithStripTitles() NormalizeOption {
	return func(c *normalizeConfig) { c.stripTitles = true }
}

// NormalizeChapters cleans up a chapter list: it drops negative start times,
// optionally drops chapters past a maximum start, deduplicates repeated start
// times, optionally strips HTML from titles, and (by default) sorts by start
// time. Deduplication keeps the first chapter seen for a given start time, and
// sorting is stable.
func NormalizeChapters(chapters []Chapter, opts ...NormalizeOption) []Chapter {
	cfg := normalizeConfig{sort: true, dedupe: true}
	for _, opt := range opts {
		opt(&cfg)
	}

	cleaned := make([]Chapter, 0, len(chapters))
	seenStarts := make(map[int]struct{})
	for _, chapter := range chapters {
		if chapter.Start < 0 {
			continue
		}
		if cfg.maxStart != nil && chapter.Start > *cfg.maxStart {
			continue
		}
		if cfg.dedupe {
			if _, ok := seenStarts[chapter.Start]; ok {
				continue
			}
		}
		seenStarts[chapter.Start] = struct{}{}
		if cfg.stripTitles {
			chapter.Title = StripHTML(chapter.Title)
		}
		cleaned = append(cleaned, chapter)
	}

	if cfg.sort {
		sort.SliceStable(cleaned, func(i, j int) bool {
			return cleaned[i].Start < cleaned[j].Start
		})
	}
	return cleaned
}
