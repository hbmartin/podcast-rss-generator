package transcript_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const specJSON = `{"version": "1.0.0", "segments": [{"startTime": 0.0, "body": "Hello world."}]}`

func TestJSONFileToJSONFileCopiesSpecFile(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "in.json")
	writeFile(t, source, specJSON)
	dest := filepath.Join(t.TempDir(), "out.json")

	require.NoError(t, transcript.JSONFileToJSONFile(source, dest, nil))

	data := decodeJSONFile(t, dest)
	assert.Equal(t, "1.0.0", data["version"])
	assert.Equal(t, "Hello world.", firstSegment(t, data)["body"])
	assert.NotContains(t, data, "metadata")
}

func TestJSONFileToJSONFileAddsMetadata(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "in.json")
	writeFile(t, source, specJSON)
	dest := filepath.Join(t.TempDir(), "out.json")

	require.NoError(t, transcript.JSONFileToJSONFile(source, dest, map[string]string{"title": "Episode 1"}))

	assert.Equal(t, map[string]any{"title": "Episode 1"}, decodeJSONFile(t, dest)["metadata"])
}

func TestJSONFileToJSONFileInvalidSchemas(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"non object":       `{"foo": "bar"}`,
		"missing segments": `{"version": "1.0.0"}`,
		"top-level array":  `[]`,
		"version not str":  `{"version": 1, "segments": []}`,
		"segments not arr": `{"version": "1.0.0", "segments": {}}`,
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			source := filepath.Join(t.TempDir(), "in.json")
			writeFile(t, source, content)
			dest := filepath.Join(t.TempDir(), "out.json")

			err := transcript.JSONFileToJSONFile(source, dest, nil)
			var jsonErr *transcript.InvalidJSONError
			require.ErrorAs(t, err, &jsonErr)
			assert.Contains(t, err.Error(), source)
			assert.NoFileExists(t, dest)
		})
	}
}

func TestJSONFileToJSONFileInvalidJSONWrapsDecodeError(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "in.json")
	writeFile(t, source, "{not valid json")
	dest := filepath.Join(t.TempDir(), "out.json")

	err := transcript.JSONFileToJSONFile(source, dest, nil)
	var jsonErr *transcript.InvalidJSONError
	require.ErrorAs(t, err, &jsonErr)
	var syntaxErr *json.SyntaxError
	require.ErrorAs(t, err, &syntaxErr)
	assert.Contains(t, err.Error(), source)
	assert.NoFileExists(t, dest)
}
