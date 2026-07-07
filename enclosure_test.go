package podcast_test

import (
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
	"github.com/stretchr/testify/assert"
)

type enclosureTest struct {
	t        podcast.EnclosureType
	expected string
}

var enclosureTests = []enclosureTest{
	{podcast.EnclosureUnknown, "application/octet-stream"},
	{podcast.M4A, "audio/x-m4a"},
	{podcast.M4V, "video/x-m4v"},
	{podcast.MP4, "video/mp4"},
	{podcast.MP3, "audio/mpeg"},
	{podcast.MOV, "video/quicktime"},
	{podcast.PDF, "application/pdf"},
	{podcast.EPUB, "document/x-epub"},
	{99, "application/octet-stream"},
}

func TestEnclosureTypes(t *testing.T) {
	t.Parallel()
	for _, et := range enclosureTests {
		t.Run(et.t.String(), func(t *testing.T) {
			t.Parallel()

			assert.EqualValues(t, et.expected, et.t.String())
		})
	}
}
