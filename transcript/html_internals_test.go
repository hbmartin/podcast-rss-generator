package transcript

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTMLTsToSecs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string
		want  int
	}{
		{"00:00:00", 0},
		{"01:31", 91},
		{"02:01:31", 7291},
	}
	for _, tt := range cases {
		got, err := htmlTsToSecs(tt.input)
		require.NoError(t, err, tt.input)
		assert.Equal(t, tt.want, got, tt.input)
	}
}

func TestHTMLTsToSecsErrors(t *testing.T) {
	t.Parallel()
	for _, input := range []string{"abc", ""} {
		_, err := htmlTsToSecs(input)
		assert.Error(t, err, input)
	}
}
