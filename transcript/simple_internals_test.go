package transcript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberToTs(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "00:00:00", numberToTs(0))
	assert.Equal(t, "00:01:01", numberToTs(61.5))
	assert.Equal(t, "01:01:01", numberToTs(3661))
}
