package transcript_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript"
	"github.com/stretchr/testify/require"
)

// TestGoldenFixtures converts every fixture with the Go converters and compares
// the parsed result to the reference JSON produced by the original Python
// package, guaranteeing byte-for-value fidelity across the whole transcript, not
// just the first and last segments.
func TestGoldenFixtures(t *testing.T) {
	t.Parallel()
	fixtures := []string{
		"300 multiple choices.html",
		"78 Exploring MCMC Sampler Algorithms, with Matt D. Hoffman.html",
		"AI Autocomplete for QGIS.srt",
		"Bloat.srt",
		"Episode 17.srt",
		"Hunting CrossSite Scripting on the Web.xsl",
		"I Quit My Job.srt",
		"Labs Setting Engineering Goals and Reporting to Stakeholders.txt",
		"Managed services vs. DIY.vtt",
		"Tales from Manufacturing Shipping Rack 1.html",
		"Talking AI at OpenShift Commons Gathering in Raleigh.html",
		"Yak Shaving with Tim Mitra.srt",
		"Zenlytic Is Building You A Better Coworker With AI Agents.txt",
	}
	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			source := filepath.Join("testdata", name)
			dest := filepath.Join(t.TempDir(), "out.json")
			require.NoError(t, transcript.ConvertFile(context.Background(), source, dest, nil))

			got := readJSON(t, dest)
			want := readJSON(t, filepath.Join("testdata", "golden", name+".json"))
			require.Equal(t, want, got, "converted output differs from Python reference for %s", name)
		})
	}
}

func readJSON(t *testing.T, path string) any {
	t.Helper()
	data, err := os.ReadFile(path) //nolint:gosec // test-controlled path
	require.NoError(t, err)
	var value any
	require.NoError(t, json.Unmarshal(data, &value))
	return value
}
