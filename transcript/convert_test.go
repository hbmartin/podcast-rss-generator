package transcript_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertFileSRT(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "ep.srt")
	writeFile(t, source, validSRT)
	dest := filepath.Join(t.TempDir(), "ep.json")

	require.NoError(t, transcript.ConvertFile(context.Background(), source, dest, nil))

	data := decodeJSONFile(t, dest)
	assert.Equal(t, "1.0.0", data["version"])
	segment := firstSegment(t, data)
	assert.Equal(t, "Michael", segment["speaker"])
	assert.Equal(t, "Hello world.", segment["body"])
}

func TestConvertFileUnknownType(t *testing.T) {
	t.Parallel()
	source := filepath.Join(t.TempDir(), "mystery.bin")
	writeFile(t, source, "not a transcript")

	err := transcript.ConvertFile(context.Background(), source, filepath.Join(t.TempDir(), "out.json"), nil)
	var unknownErr *transcript.UnknownFileTypeError
	require.ErrorAs(t, err, &unknownErr)
}

func TestBulkConvert(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	source := filepath.Join(dir, "in")
	require.NoError(t, os.MkdirAll(filepath.Join(source, "show"), 0o750))
	writeFile(t, filepath.Join(source, "show", "ep1.srt"), validSRT)
	writeFile(t, filepath.Join(source, "show", "ep2.vtt"), validVTT)
	writeFile(t, filepath.Join(source, "mystery.bin"), "not a transcript")
	dest := filepath.Join(dir, "out")

	summary, err := transcript.BulkConvert(context.Background(), source, dest)
	require.NoError(t, err)
	assert.Len(t, summary.Converted, 2)
	assert.Empty(t, summary.Failed)
	assert.Equal(t, []string{filepath.Join(source, "mystery.bin")}, summary.Unknown)
	assert.FileExists(t, filepath.Join(dest, "show", "ep1.json"))
	assert.FileExists(t, filepath.Join(dest, "show", "ep2.json"))
}

func TestBulkConvertContinuesAfterFailure(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	source := filepath.Join(dir, "in")
	require.NoError(t, os.Mkdir(source, 0o750))
	writeFile(t, filepath.Join(source, "bad.srt"), "this is not valid srt")
	writeFile(t, filepath.Join(source, "good.srt"), validSRT)
	dest := filepath.Join(dir, "out")

	summary, err := transcript.BulkConvert(context.Background(), source, dest)
	require.NoError(t, err)
	assert.Equal(t, []string{filepath.Join(source, "good.srt")}, convertedSources(summary))
	assert.Equal(t, []string{filepath.Join(source, "bad.srt")}, failedSources(summary))
	assert.FileExists(t, filepath.Join(dest, "good.json"))
}

func TestBulkConvertRecordsInvalidJSONAsFailure(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	source := filepath.Join(dir, "in")
	require.NoError(t, os.Mkdir(source, 0o750))
	invalidJSON := filepath.Join(source, "bad.json")
	writeFile(t, invalidJSON, `{"version": "1.0.0"}`)
	dest := filepath.Join(dir, "out")

	summary, err := transcript.BulkConvert(context.Background(), source, dest)
	require.NoError(t, err)
	assert.Empty(t, summary.Converted)
	assert.Equal(t, []string{invalidJSON}, failedSources(summary))
	assert.NoFileExists(t, filepath.Join(dest, "bad.json"))
}

func TestBulkConvertRecordsMissingVTTAsFailure(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.vtt")
	dbPath := filepath.Join(dir, "overcast.db")
	dest := filepath.Join(dir, "out")

	lister := &capturingDBLister{paths: []string{missing}}
	summary, err := transcript.BulkConvert(context.Background(), dbPath, dest, transcript.WithDBLister(lister))
	require.NoError(t, err)
	assert.Equal(t, dbPath, lister.gotDBPath)
	assert.Empty(t, lister.gotIgnore)
	assert.Empty(t, summary.Converted)
	assert.Equal(t, []string{missing}, failedSources(summary))
	assert.Empty(t, jsonFilesUnder(t, dest))
}

// jsonFilesUnder returns every .json file under dir (recursively).
func jsonFilesUnder(t *testing.T, dir string) []string {
	t.Helper()
	var found []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr == nil && !d.IsDir() && filepath.Ext(path) == ".json" {
			found = append(found, path)
		}
		return nil
	})
	require.NoError(t, err)
	return found
}

func TestBulkConvertSkipsExistingUnlessOverwrite(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	source := filepath.Join(dir, "in")
	require.NoError(t, os.Mkdir(source, 0o750))
	writeFile(t, filepath.Join(source, "ep.srt"), validSRT)
	dest := filepath.Join(dir, "out")

	first, err := transcript.BulkConvert(context.Background(), source, dest)
	require.NoError(t, err)
	assert.Len(t, first.Converted, 1)

	second, err := transcript.BulkConvert(context.Background(), source, dest)
	require.NoError(t, err)
	assert.Empty(t, second.Converted)
	assert.Len(t, second.Skipped, 1)

	third, err := transcript.BulkConvert(context.Background(), source, dest, transcript.WithOverwrite())
	require.NoError(t, err)
	assert.Len(t, third.Converted, 1)
	assert.Empty(t, third.Skipped)
}

func TestBulkConvertDryRunWritesNothing(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	source := filepath.Join(dir, "in")
	require.NoError(t, os.Mkdir(source, 0o750))
	writeFile(t, filepath.Join(source, "ep.srt"), validSRT)
	dest := filepath.Join(dir, "out")

	summary, err := transcript.BulkConvert(context.Background(), source, dest, transcript.WithDryRun())
	require.NoError(t, err)
	assert.True(t, summary.DryRun)
	assert.Len(t, summary.Converted, 1)
	assert.NoDirExists(t, dest)
}

func TestBulkConvertIgnoreList(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	source := filepath.Join(dir, "in")
	require.NoError(t, os.Mkdir(source, 0o750))
	writeFile(t, filepath.Join(source, "ep.srt"), validSRT)
	writeFile(t, filepath.Join(source, "skipme.srt"), validSRT)
	dest := filepath.Join(dir, "out")

	summary, err := transcript.BulkConvert(context.Background(), source, dest,
		transcript.WithIgnore([]string{"skipme.srt"}))
	require.NoError(t, err)
	assert.Equal(t, []string{filepath.Join(source, "ep.srt")}, convertedSources(summary))
}

// capturingDBLister is a fake transcript.DBLister for BulkConvert tests.
type capturingDBLister struct {
	paths     []string
	metadata  map[string]map[string]string
	gotDBPath string
	gotIgnore []string
}

func (l *capturingDBLister) ListFiles(
	_ context.Context,
	dbPath string,
	ignore []string,
) ([]string, map[string]map[string]string, error) {
	l.gotDBPath = dbPath
	l.gotIgnore = ignore
	return l.paths, l.metadata, nil
}

func convertedSources(summary *transcript.ConversionSummary) []string {
	sources := make([]string, len(summary.Converted))
	for i, pair := range summary.Converted {
		sources[i] = pair.Source
	}
	return sources
}

func failedSources(summary *transcript.ConversionSummary) []string {
	sources := make([]string, len(summary.Failed))
	for i, failure := range summary.Failed {
		sources[i] = failure.Source
	}
	return sources
}
