package transcript

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
)

// SpecVersion is the PodcastIndex transcript JSON spec version emitted by this
// package.
//
// https://github.com/Podcastindex-org/podcast-namespace/blob/main/transcripts/transcripts.md
const SpecVersion = "1.0.0"

// speakerPrefixRe matches a leading "Speaker Name:" prefix, e.g. "Michael:
// Hello" or "Leo Dion (host): Before we begin". The name must start with an
// uppercase letter, be reasonably short, and be followed by a colon and
// whitespace.
var speakerPrefixRe = regexp.MustCompile(`(?s)^([A-Z][A-Za-z0-9 .\-'()&]{0,62}):\s+(\S.*)$`)

// SplitSpeakerPrefix splits a leading "Speaker Name:" prefix from a segment
// body, returning (speaker, body). The speaker is empty when no prefix is
// found.
func SplitSpeakerPrefix(body string) (speaker, rest string) {
	if match := speakerPrefixRe.FindStringSubmatch(body); match != nil {
		return strings.TrimSpace(match[1]), match[2]
	}
	return "", body
}

// Segment is a single transcript segment (one cue or utterance). A nil field is
// omitted from the serialized output, distinguishing an unset field from an
// empty string or a zero time.
type Segment struct {
	Speaker   *string  `json:"speaker,omitempty"`
	StartTime *float64 `json:"startTime,omitempty"`
	EndTime   *float64 `json:"endTime,omitempty"`
	Body      *string  `json:"body,omitempty"`
}

// isEmpty reports whether every field is unset (the segment would serialize to
// an empty object).
func (s Segment) isEmpty() bool {
	return s.Speaker == nil && s.StartTime == nil && s.EndTime == nil && s.Body == nil
}

// Transcript is a full transcript: a list of segments plus optional feed
// metadata.
type Transcript struct {
	Version  string            `json:"version"`
	Segments []Segment         `json:"segments"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// newTranscript returns a Transcript with the given segments and the default
// spec version.
func newTranscript(segments []Segment) *Transcript {
	return &Transcript{Version: SpecVersion, Segments: segments}
}

func strPtr(s string) *string     { return &s }
func floatPtr(f float64) *float64 { return &f }

// encodeJSON serializes v as UTF-8 JSON without HTML escaping (matching
// Python's json.dumps), indented with the given number of spaces, and without a
// trailing newline.
func encodeJSON(v any, indent int) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if indent > 0 {
		encoder.SetIndent("", strings.Repeat(" ", indent))
	}
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
