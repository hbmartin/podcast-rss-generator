package transcript

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMtsToSecsFloatProperty(t *testing.T) {
	t.Parallel()
	for _, hours := range []int{0, 1, 23, 99} {
		for _, minutes := range []int{0, 30, 59} {
			for _, seconds := range []int{0, 30, 59} {
				for _, millis := range []int{0, 1, 123, 999} {
					want := math.Round((float64(hours*3600+minutes*60+seconds)+float64(millis)/1000)*1000) / 1000
					comma := fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, millis)
					period := fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)

					gotComma, err := mtsToSecsFloat(comma)
					require.NoError(t, err)
					assert.InDelta(t, want, gotComma, 0, comma)

					gotPeriod, err := mtsToSecsFloat(period)
					require.NoError(t, err)
					assert.InDelta(t, want, gotPeriod, 0, period)
				}
			}
		}
	}
}

func TestHTMLTsToSecsProperty(t *testing.T) {
	t.Parallel()
	for _, hours := range []int{0, 1, 12, 23} {
		for _, minutes := range []int{0, 30, 59} {
			for _, seconds := range []int{0, 30, 59} {
				mmss := fmt.Sprintf("%d:%02d", minutes, seconds)
				got, err := htmlTsToSecs(mmss)
				require.NoError(t, err)
				assert.Equal(t, minutes*60+seconds, got, mmss)

				hhmmss := fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
				got, err = htmlTsToSecs(hhmmss)
				require.NoError(t, err)
				assert.Equal(t, hours*3600+minutes*60+seconds, got, hhmmss)
			}
		}
	}
}

func TestSplitSpeakerPrefixNeverCrashes(t *testing.T) {
	t.Parallel()
	property := func(body string) bool {
		speaker, remainder := SplitSpeakerPrefix(body)
		if !strings.HasSuffix(body, remainder) {
			return false
		}
		if speaker != "" {
			if !strings.HasPrefix(body, speaker) {
				return false
			}
			if speaker != strings.TrimSpace(speaker) {
				return false
			}
		}
		return true
	}
	require.NoError(t, quick.Check(property, &quick.Config{MaxCount: 2000}))
}
