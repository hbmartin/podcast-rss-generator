package transcript_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFileToSimpleFile(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "ep.json")
	writeFile(t, source, `{"version":"1.0.0","segments":[`+
		`{"startTime":0.0,"body":"Hello.","speaker":"Ann"},`+
		`{"startTime":61.0,"body":"Goodbye."}]}`)
	dest := filepath.Join(t.TempDir(), "ep.txt")

	require.NoError(t, transcript.JSONFileToSimpleFile(source, dest))

	data, err := os.ReadFile(dest) //nolint:gosec // test-controlled path
	require.NoError(t, err)
	assert.Equal(t, "(00:00:00) Ann: Hello.\n(00:01:01) Goodbye.", string(data))
}

func TestJSONFileToSimpleFileWithoutStartTime(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "ep.json")
	writeFile(t, source, `{"segments":[{"body":"Hello."}]}`)

	err := transcript.JSONFileToSimpleFile(source, filepath.Join(t.TempDir(), "ep.txt"))
	var noStart *transcript.NoStartTimeError
	require.ErrorAs(t, err, &noStart)
	assert.Contains(t, err.Error(), source)
}

func TestJSONFileToSimpleFileInvalidJSON(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "ep.json")
	writeFile(t, source, "{not valid json")

	err := transcript.JSONFileToSimpleFile(source, filepath.Join(t.TempDir(), "ep.txt"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), source)
}
