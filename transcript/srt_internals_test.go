package transcript

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMtsToSecsFloat(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string
		want  float64
	}{
		{"00:00:00,0", 0},
		{"00:01:31,123", 91.123},
		{"23:59:59,999", 86399.999},
		{"02:01:31,0", 7291},
	}
	for _, tt := range cases {
		got, err := mtsToSecsFloat(tt.input)
		require.NoError(t, err, tt.input)
		assert.InDelta(t, tt.want, got, 0, tt.input)
	}
}

func TestMtsToSecsFloatErrors(t *testing.T) {
	t.Parallel()
	for _, input := range []string{"abc", ""} {
		_, err := mtsToSecsFloat(input)
		assert.Error(t, err, input)
	}
}
