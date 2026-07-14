// Package chapters extracts and transforms podcast chapter information between
// PodcastIndex (PCI) chapters, Podlove Simple Chapters (PSC), description
// embeds (timestamps in show notes), and ID3v2 CHAP frames.
//
// It is a Go port of the Python package podcast-chapter-tools. Every extractor
// returns a slice of [Chapter] (nil when nothing usable is found), where the
// start time is expressed in whole seconds.
package chapters

import (
	"regexp"

	"github.com/dlclark/regexp2"
)

// Namespaces used by the PodcastIndex and Podlove Simple Chapters formats.
const (
	pciNamespace = "https://podcastindex.org/namespace/1.0"
	pscNamespace = "http://podlove.org/simple-chapters"
)

// Chapter is a single podcast chapter. Start is the chapter start time in whole
// seconds. URL and Image are empty strings when absent, mirroring the optional
// url/image fields of the source formats.
type Chapter struct {
	Start int
	Title string
	URL   string
	Image string
}

// ChapterType enumerates the recognized chapter provenance labels. The AI value
// supersedes the deprecated Gemini, OpenAI, and Sonnet labels.
type ChapterType int

// ChapterType values.
const (
	ChapterTypeAI ChapterType = iota
	ChapterTypeDescription
	ChapterTypeGemini // Deprecated: use ChapterTypeAI.
	ChapterTypeID3
	ChapterTypeOpenAI // Deprecated: use ChapterTypeAI.
	ChapterTypePCI
	ChapterTypePSC
	ChapterTypeSonnet // Deprecated: use ChapterTypeAI.
)

// String returns the lowercase label of the ChapterType.
func (t ChapterType) String() string {
	switch t {
	case ChapterTypeAI:
		return "ai"
	case ChapterTypeDescription:
		return "description"
	case ChapterTypeGemini:
		return "gemini"
	case ChapterTypeID3:
		return "id3"
	case ChapterTypeOpenAI:
		return "openai"
	case ChapterTypePCI:
		return "pci"
	case ChapterTypePSC:
		return "psc"
	case ChapterTypeSonnet:
		return "sonnet"
	}
	return ""
}

// Regexes for extracting chapters from an episode description. The two chapter
// patterns require backtracking features (lazy quantifiers and, for the retry
// pattern, a lookahead) that Go's RE2 engine does not provide, so they use the
// regexp2 backtracking engine to match the original Python semantics exactly.
var (
	descriptionChapterRe = regexp2.MustCompile(
		`>?\s?[(\[]?(\d{0,2}:?\d{1,2}:\d{2})[\])]?[\s-]*([^\[(\n]+?)\s*(?:<|$)`,
		regexp2.Multiline,
	)
	retryDescriptionChapterRe = regexp2.MustCompile(
		`(\d{0,2}:?\d{1,2}:\d{2})[\])]?[\s-]*([^\[(]+?)(?=\d{1,2}:|$)`,
		regexp2.Multiline,
	)

	urlRe     = regexp.MustCompile(`(?i)https?://[^\s"'<>]+`)
	htmlTagRe = regexp.MustCompile(`<[^>]+>`)
)

// findAllPairs returns the (group 1, group 2) captures of every match of re in
// s, mirroring Python's re.findall on a two-group pattern.
func findAllPairs(re *regexp2.Regexp, s string) [][2]string {
	var out [][2]string
	m, err := re.FindStringMatch(s)
	for err == nil && m != nil {
		groups := m.Groups()
		out = append(out, [2]string{groups[1].String(), groups[2].String()})
		m, err = re.FindNextMatch(m)
	}
	return out
}
