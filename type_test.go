package podcast_test

import (
	"testing"

	"github.com/eduncan911/podcast"
	"github.com/stretchr/testify/assert"
)

func TestPodcastTypeStringSpecValues(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "episodic", podcast.Episodic.String())
	assert.Equal(t, "serial", podcast.Serial.String())
	assert.Equal(t, "episodic", podcast.PodcastType(99).String())
}

func TestEpisodeTypeStringSpecValues(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "full", podcast.Full.String())
	assert.Equal(t, "trailer", podcast.Trailer.String())
	assert.Equal(t, "bonus", podcast.Bonus.String())
	assert.Equal(t, "full", podcast.EpisodeType(99).String())
}

func TestAddTypeUsesSpecValue(t *testing.T) {
	t.Parallel()

	p := podcast.New("title", "link", "description", zeroDate, zeroDate)

	p.AddType(podcast.Serial)

	if assert.NotNil(t, p.IType) {
		assert.Equal(t, "serial", p.IType.Text)
	}
	assert.Contains(t, p.String(), "<itunes:type>serial</itunes:type>")
}

func TestAddChannelTypeUsesSpecValue(t *testing.T) {
	t.Parallel()

	p := podcast.New("title", "link", "description", zeroDate, zeroDate)

	p.AddChannelType(" Serial ")

	if assert.NotNil(t, p.IType) {
		assert.Equal(t, "serial", p.IType.Text)
	}
}
