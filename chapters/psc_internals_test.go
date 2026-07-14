package chapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIterFeedItemsWithoutChannel(t *testing.T) {
	t.Parallel()
	root, err := ParseXML([]byte("<rss></rss>"))
	require.NoError(t, err)
	assert.Empty(t, iterFeedItems(root))
}
