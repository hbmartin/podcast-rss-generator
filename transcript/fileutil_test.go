package transcript_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFiles(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "sub"), 0o750))
	writeFile(t, filepath.Join(dir, "a.srt"), "x")
	writeFile(t, filepath.Join(dir, "sub", "b.vtt"), "x")
	writeFile(t, filepath.Join(dir, ".hidden"), "x")
	writeFile(t, filepath.Join(dir, "skip.pdf"), "x")
	writeFile(t, filepath.Join(dir, "ignored.srt"), "x")

	files, err := transcript.ListFiles(dir, []string{"ignored.srt"})
	require.NoError(t, err)
	sort.Strings(files)
	assert.Equal(t, []string{filepath.Join(dir, "a.srt"), filepath.Join(dir, "sub", "b.vtt")}, files)
}

func TestReadTextRobustUTF8BOM(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "bom.srt")
	require.NoError(t, os.WriteFile(path, []byte("\xef\xbb\xbfWEBVTT\n"), 0o600))
	got, err := transcript.ReadTextRobust(path)
	require.NoError(t, err)
	assert.Equal(t, "WEBVTT\n", got)
}

func TestReadTextRobustInvalidUTF8(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "latin1.srt")
	require.NoError(t, os.WriteFile(path, []byte("caf\xe9 time\n"), 0o600))
	got, err := transcript.ReadTextRobust(path)
	require.NoError(t, err)
	assert.Equal(t, "caf� time\n", got)
}

func TestReadFirstLineInvalidUTF8(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "latin1.srt")
	require.NoError(t, os.WriteFile(path, []byte("caf\xe9 time\nsecond line\n"), 0o600))
	got, err := transcript.ReadFirstLine(path)
	require.NoError(t, err)
	assert.Equal(t, "caf� time\n", got)
}

func TestMapFilesInParallelPreservesOrder(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	const fileCount = 5
	paths := make([]string, 0, fileCount)
	for i := range fileCount {
		path := filepath.Join(dir, strconv.Itoa(i)+".txt")
		writeFile(t, path, strconv.Itoa(i))
		paths = append(paths, path)
	}

	results := transcript.MapFilesInParallel(context.Background(), paths, func(_ context.Context, p string) string {
		data, err := os.ReadFile(p) //nolint:gosec // test-controlled path
		if err != nil {
			return ""
		}
		return string(data)
	})

	require.Len(t, results, len(paths))
	for i, path := range paths {
		assert.Equal(t, path, results[i].Path)
		assert.Equal(t, strconv.Itoa(i), results[i].Result)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
}
