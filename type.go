package podcast

// Episodic and Serial are the supported PodcastType values.
const (
	Episodic PodcastType = iota
	Serial
)

const (
	podcastTypeDefault = "episodic"
)

// PodcastType specifies the type of the podcast.
//
// Its values can be one of the following:
//
// Episodic (default). Specify episodic when episodes are intended to be
// consumed without any specific order. Apple Podcasts will present newest
// episodes first and display the publish date (required) of each episode.
// If organized into seasons, the newest season will be presented first
// - otherwise, episodes will be grouped by year published, newest first.
//
// For new subscribers, Apple Podcasts adds the newest, most recent episode
// in their Library.
//
// Serial. Specify serial when episodes are intended to be consumed in
// sequential order. Apple Podcasts will present the oldest episodes
// first and display the episode numbers (required) of each episode. If
// organized into seasons, the newest season will be presented first and
// <itunes:episode> numbers must be given for each episode.
//
// For new subscribers, Apple Podcasts adds the first episode to their
// Library, or the entire current season if using seasons.
//
//nolint:revive // Preserve the established public API name.
type PodcastType int

// String returns the Apple Podcasts encoding of the specified PodcastType.
func (pt PodcastType) String() string {
	// https://help.apple.com/itc/podcasts_connect/#/itcb54353390
	switch pt {
	case Episodic:
		return "episodic"
	case Serial:
		return "serial"
	}
	return podcastTypeDefault
}

// EpisodeType specifies the type of an episode.
const (
	Full EpisodeType = iota
	Trailer
	Bonus
)

const (
	episodeTypeDefault = "full"
)

// EpisodeType specifies whether an episode is full, trailer, or bonus content.
type EpisodeType int

// String returns the Apple Podcasts encoding of the specified EpisodeType.
func (et EpisodeType) String() string {
	switch et {
	case Full:
		return "full"
	case Trailer:
		return "trailer"
	case Bonus:
		return "bonus"
	}
	return episodeTypeDefault
}
