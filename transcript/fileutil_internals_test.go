package transcript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFileAllowed(t *testing.T) {
	t.Parallel()
	assert.True(t, isFileAllowed("a.srt", nil))
	assert.False(t, isFileAllowed("a.srt", []string{"a.srt"}))
	assert.False(t, isFileAllowed(".hidden", nil))
	assert.False(t, isFileAllowed("doc.pdf", nil))
	assert.False(t, isFileAllowed("blob.octet-stream", nil))
}
