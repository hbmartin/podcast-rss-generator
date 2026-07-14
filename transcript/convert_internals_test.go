package transcript

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validSRTContent = "1\n00:00:00,000 --> 00:00:01,000\nMichael: Hello world.\n\n"

func TestDestinationPathMirrorsSourceRoot(t *testing.T) {
	t.Parallel()
	assert.Equal(t, filepath.Join("/out", "show", "ep.json"),
		destinationPath("/src/show/ep.srt", "/out", "/src"))
}

func TestDestinationPathNestedDirsDoNotCollide(t *testing.T) {
	t.Parallel()
	a := destinationPath("/src/show a/deep/ep.srt", "/out", "/src")
	b := destinationPath("/src/show b/deep/ep.srt", "/out", "/src")
	assert.NotEqual(t, a, b)
}

func TestDestinationPathWithoutSourceRootUsesParentName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, filepath.Join("/out", "show", "ep.json"),
		destinationPath("/anywhere/show/ep.srt", "/out", ""))
}

func TestDestinationPathSourceRootNotAncestorFallsBack(t *testing.T) {
	t.Parallel()
	assert.Equal(t, filepath.Join("/out", "show", "ep.json"),
		destinationPath("/other/show/ep.srt", "/out", "/src"))
}

func TestDestinationPathBareFilename(t *testing.T) {
	t.Parallel()
	assert.Equal(t, filepath.Join("/out", "ep.json"), destinationPath("ep.srt", "/out", ""))
	assert.Equal(t, filepath.Join("/out", "ep.json"), destinationPath("/ep.srt", "/out", ""))
}

func TestAssignDestinationsDedupesCollisions(t *testing.T) {
	t.Parallel()
	destinations := assignDestinations([]string{"/db/show/ep.srt", "/db/show/ep.vtt"}, "/out", "")
	assert.Equal(t, filepath.Join("/out", "show", "ep.json"), destinations["/db/show/ep.srt"])
	assert.Equal(t, filepath.Join("/out", "show", "ep (1).json"), destinations["/db/show/ep.vtt"])
}

func TestBulkConvertDuplicateSourcesDoNotCrash(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	source := filepath.Join(dir, "in")
	require.NoError(t, os.Mkdir(source, 0o750))
	transcriptFile := filepath.Join(source, "ep.srt")
	require.NoError(t, os.WriteFile(transcriptFile, []byte(validSRTContent), 0o600))
	dest := filepath.Join(dir, "out")

	duplicate := func(directory string, ignore []string) ([]string, error) {
		assert.Equal(t, source, directory)
		assert.Empty(t, ignore)
		return []string{transcriptFile, transcriptFile}, nil
	}

	summary, err := BulkConvert(context.Background(), source, dest, WithDryRun(), withFileLister(duplicate))
	require.NoError(t, err)
	sources := make([]string, len(summary.Converted))
	for i, pair := range summary.Converted {
		sources[i] = pair.Source
	}
	assert.Equal(t, []string{transcriptFile, transcriptFile}, sources)
}
