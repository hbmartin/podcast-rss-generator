package chapters_test

import (
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTsToSecs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"seconds only", "42", 42},
		{"minutes seconds", "1:02", 62},
		{"hours minutes seconds", "1:02:03", 3723},
		{"fractional seconds truncated", "0:01:30.500", 90},
		{"leading zeros zero", "00:00:00", 0},
		{"leading zeros minutes", "01:05", 65},
		{"surrounding whitespace", " 1:02 ", 62},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := chapters.TsToSecs(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTsToSecsErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantMsg string
	}{
		{"fractional not in final segment", "1.5:02", "final segment"},
		{"fractional not numeric", "1:02.bad", ""},
		{"too many segments", "1:2:3:4", ""},
		{"non numeric two parts", "aa:10", ""},
		{"non numeric three parts", "aa:10:10", ""},
		{"out of range minutes", "1:61:00", ""},
		{"out of range seconds", "1:00:99", ""},
		{"negative segment", "-5", "negative segment"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := chapters.TsToSecs(tt.input)
			require.Error(t, err)
			if tt.wantMsg != "" {
				assert.Contains(t, err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestSecsToTs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		seconds int
		want    string
	}{
		{"under a minute", 5, "0:05"},
		{"minutes", 62, "1:02"},
		{"hours", 3723, "1:02:03"},
		{"zero", 0, "0:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := chapters.SecsToTs(tt.seconds)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecsToTsNegative(t *testing.T) {
	t.Parallel()
	_, err := chapters.SecsToTs(-1)
	require.Error(t, err)
}

func TestTimecodesRoundtrip(t *testing.T) {
	t.Parallel()
	for _, seconds := range []int{0, 59, 60, 61, 3599, 3600, 3661, 86399} {
		ts, err := chapters.SecsToTs(seconds)
		require.NoError(t, err)
		got, err := chapters.TsToSecs(ts)
		require.NoError(t, err)
		assert.Equal(t, seconds, got)
	}
}
