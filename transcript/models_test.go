package transcript_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strptr(s string) *string   { return &s }
func f64ptr(f float64) *float64 { return &f }

// toMap serializes v and decodes it into a generic map for comparison.
func toMap(t *testing.T, v any) map[string]any {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)
	var m map[string]any
	require.NoError(t, json.Unmarshal(data, &m))
	return m
}

func TestSplitSpeakerPrefix(t *testing.T) {
	t.Parallel()
	speaker, body := transcript.SplitSpeakerPrefix("Michael: Hello there")
	assert.Equal(t, "Michael", speaker)
	assert.Equal(t, "Hello there", body)

	speaker, body = transcript.SplitSpeakerPrefix("Leo Dion (host): Before we begin")
	assert.Equal(t, "Leo Dion (host)", speaker)
	assert.Equal(t, "Before we begin", body)

	speaker, body = transcript.SplitSpeakerPrefix("Dr. Smith: Hi")
	assert.Equal(t, "Dr. Smith", speaker)
	assert.Equal(t, "Hi", body)
}

func TestSplitSpeakerPrefixNoMatch(t *testing.T) {
	t.Parallel()
	cases := []string{
		"Hello there",
		"lowercase: nope",
		"Michael:no space",
		"",
		strings.Repeat("A", 80) + ": too long",
	}
	for _, input := range cases {
		speaker, body := transcript.SplitSpeakerPrefix(input)
		assert.Empty(t, speaker)
		assert.Equal(t, input, body)
	}
}

func TestSegmentSerializationOmitsUnsetFields(t *testing.T) {
	t.Parallel()
	assert.Equal(t, map[string]any{}, toMap(t, transcript.Segment{}))
	assert.Equal(t, map[string]any{"body": "hi"}, toMap(t, transcript.Segment{Body: strptr("hi")}))
	assert.Equal(t, map[string]any{"startTime": 0.0, "body": ""},
		toMap(t, transcript.Segment{Body: strptr(""), StartTime: f64ptr(0.0)}))
	assert.Equal(t, map[string]any{"speaker": "Ann", "startTime": 1.0, "endTime": 2.0, "body": "hi"},
		toMap(t, transcript.Segment{
			Body:      strptr("hi"),
			StartTime: f64ptr(1.0),
			EndTime:   f64ptr(2.0),
			Speaker:   strptr("Ann"),
		}))
}

func TestSegmentKeyOrder(t *testing.T) {
	t.Parallel()
	data, err := json.Marshal(transcript.Segment{
		Speaker:   strptr("Ann"),
		StartTime: f64ptr(1),
		EndTime:   f64ptr(2),
		Body:      strptr("hi"),
	})
	require.NoError(t, err)
	const want = `{"speaker":"Ann","startTime":1,"endTime":2,"body":"hi"}`
	if got := string(data); got != want {
		t.Errorf("segment key order:\n got %s\nwant %s", got, want)
	}
}

func TestTranscriptSerialization(t *testing.T) {
	t.Parallel()
	tr := transcript.Transcript{Version: transcript.SpecVersion, Segments: []transcript.Segment{{Body: strptr("hi")}}}
	assert.Equal(t, map[string]any{
		"version":  "1.0.0",
		"segments": []any{map[string]any{"body": "hi"}},
	}, toMap(t, tr))
}

func TestTranscriptSerializationWithMetadata(t *testing.T) {
	t.Parallel()
	tr := transcript.Transcript{
		Version:  transcript.SpecVersion,
		Segments: []transcript.Segment{{Body: strptr("hi")}},
		Metadata: map[string]string{"title": "Episode 1"},
	}
	assert.Equal(t, map[string]any{"title": "Episode 1"}, toMap(t, tr)["metadata"])
}
