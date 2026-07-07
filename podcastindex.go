package podcast

import (
	"crypto/sha1" //nolint:gosec // UUIDv5 requires SHA-1 per RFC 4122; not used for security.
	"encoding/xml"
	"fmt"
	"strings"
)

// podcastIndexNS is the Podcasting 2.0 namespace declared on <rss> as
// xmlns:podcast whenever any podcast:* tag is present in the feed.
//
// https://github.com/Podcastindex-org/podcast-namespace/blob/main/docs/1.0.md
const podcastIndexNS = "https://podcastindex.org/namespace/1.0"

const personNameRuneLimit = 128

// feedGUIDNamespace is the UUID namespace used to derive podcast:guid values,
// ead4c236-bf58-58c6-a2c6-a6b28d128cb6 per the Podcasting 2.0 specification.
var feedGUIDNamespace = [16]byte{
	0xea, 0xd4, 0xc2, 0x36, 0xbf, 0x58, 0x58, 0xc6,
	0xa2, 0xc6, 0xa6, 0xb2, 0x8d, 0x12, 0x8c, 0xb6,
}

// RFC 4122 version and variant bits applied to a UUIDv5's octets 6 and 8.
const (
	uuidVersionMask = 0x0f
	uuidVersion5    = 0x50
	uuidVariantMask = 0x3f
	uuidVariantRFC  = 0x80
)

// PPerson represents a podcast:person tag crediting a host, guest, or other
// contributor of the show or episode.
//
// Role and Group values come from the Podcast Taxonomy Project. Per the
// specification an absent role means "host" and an absent group means "cast",
// so empty attributes are omitted from the output.
//
// Use Podcast.AddPerson or Item.AddPerson to populate this correctly.
type PPerson struct {
	XMLName xml.Name `xml:"podcast:person"`
	Role    string   `xml:"role,attr,omitempty"`
	Group   string   `xml:"group,attr,omitempty"`
	Img     string   `xml:"img,attr,omitempty"`
	Href    string   `xml:"href,attr,omitempty"`
	Name    string   `xml:",chardata"`
}

// PTranscript represents a podcast:transcript tag linking to an episode
// transcript or closed-caption file.
//
// Use Item.AddTranscript to populate this correctly.
type PTranscript struct {
	XMLName  xml.Name `xml:"podcast:transcript"`
	URL      string   `xml:"url,attr"`
	Type     string   `xml:"type,attr"`
	Language string   `xml:"language,attr,omitempty"`
	Rel      string   `xml:"rel,attr,omitempty"`
}

// PChapters represents a podcast:chapters tag linking to an episode chapters
// file, typically of type "application/json+chapters".
//
// Use Item.AddChapters to populate this correctly.
type PChapters struct {
	XMLName xml.Name `xml:"podcast:chapters"`
	URL     string   `xml:"url,attr"`
	Type    string   `xml:"type,attr"`
}

// PSocialInteract represents a podcast:socialInteract tag pointing at the
// root social-media post for an episode, where comments and discussion take
// place.
//
// Use Item.AddSocialInteract to populate this correctly.
type PSocialInteract struct {
	XMLName   xml.Name `xml:"podcast:socialInteract"`
	URI       string   `xml:"uri,attr"`
	Protocol  string   `xml:"protocol,attr"`
	AccountID string   `xml:"accountId,attr,omitempty"`
}

// PLocked represents the podcast:locked tag, which tells other podcast
// platforms whether they are allowed to import this feed.
//
// Use Podcast.SetLocked to populate this correctly.
type PLocked struct {
	XMLName xml.Name `xml:"podcast:locked"`
	Owner   string   `xml:"owner,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

// MediumPodcast through MediumBlogList are the supported Medium values.
const (
	MediumPodcast Medium = iota
	MediumMusic
	MediumVideo
	MediumFilm
	MediumAudiobook
	MediumNewsletter
	MediumBlog
	MediumPublisher
	MediumCourse
	MediumMixed
	MediumPodcastList
	MediumMusicList
	MediumVideoList
	MediumFilmList
	MediumAudiobookList
	MediumNewsletterList
	MediumBlogList
)

const mediumDefault = "podcast"

// Medium specifies the podcast:medium value describing what a feed contains,
// so applications can adapt their behavior (for example resetting playback
// speed for music feeds).
//
// The *List variants encode the "L" suffix (for example "podcastL") marking
// playlist-style feeds of remote items; MediumMixed marks a feed mixing
// several remote item types.
type Medium int

// String returns the podcast namespace encoding of the specified Medium.
func (m Medium) String() string {
	switch m {
	case MediumPodcast:
		return "podcast"
	case MediumMusic:
		return "music"
	case MediumVideo:
		return "video"
	case MediumFilm:
		return "film"
	case MediumAudiobook:
		return "audiobook"
	case MediumNewsletter:
		return "newsletter"
	case MediumBlog:
		return "blog"
	case MediumPublisher:
		return "publisher"
	case MediumCourse:
		return "course"
	case MediumMixed:
		return "mixed"
	case MediumPodcastList:
		return "podcastL"
	case MediumMusicList:
		return "musicL"
	case MediumVideoList:
		return "videoL"
	case MediumFilmList:
		return "filmL"
	case MediumAudiobookList:
		return "audiobookL"
	case MediumNewsletterList:
		return "newsletterL"
	case MediumBlogList:
		return "blogL"
	}
	return mediumDefault
}

// SetPodcastGUID sets the podcast:guid tag, the globally unique identifier
// that follows a podcast throughout its lifetime across hosting platforms.
//
// The value should be a UUIDv5 derived from the feed URL; use NewFeedGUID to
// compute it. If guid is empty this method is a no-op.
func (p *Podcast) SetPodcastGUID(guid string) {
	if len(guid) == 0 {
		return
	}
	p.PGUID = guid
}

// NewFeedGUID returns the podcast:guid UUIDv5 for feedURL per the Podcasting
// 2.0 specification: the URL scheme ("https://" or "http://") and any
// trailing slashes are stripped, then the result is hashed within the podcast
// namespace UUID ead4c236-bf58-58c6-a2c6-a6b28d128cb6.
//
// If feedURL is empty after stripping, it returns "".
func NewFeedGUID(feedURL string) string {
	name := strings.TrimPrefix(feedURL, "https://")
	name = strings.TrimPrefix(name, "http://")
	name = strings.TrimRight(name, "/")
	if len(name) == 0 {
		return ""
	}
	//nolint:gosec // UUIDv5 requires SHA-1 per RFC 4122; not used for security.
	sum := sha1.Sum(append(feedGUIDNamespace[:], name...))
	sum[6] = sum[6]&uuidVersionMask | uuidVersion5
	sum[8] = sum[8]&uuidVariantMask | uuidVariantRFC
	return fmt.Sprintf("%x-%x-%x-%x-%x", sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}

// SetMedium sets the podcast:medium tag from the typed Medium value.
func (p *Podcast) SetMedium(medium Medium) {
	p.PMedium = medium.String()
}

// SetLocked sets the podcast:locked tag. Both values are meaningful and
// serialized: true renders "yes" (other platforms must not import this feed),
// false renders "no".
//
// ownerEmail is an optional address that can be used to verify feed ownership
// during move and import operations; it is rendered as the owner attribute
// when non-empty.
func (p *Podcast) SetLocked(locked bool, ownerEmail string) {
	value := "no"
	if locked {
		value = "yes"
	}
	p.PLocked = &PLocked{
		Owner: ownerEmail,
		Value: value,
	}
}

// AddPerson appends a podcast:person credit to the channel, describing a
// person of interest to the whole podcast such as a host or regular co-host.
//
// role, group, img and href are optional; empty role and group imply "host"
// and "cast" per the specification and are omitted from the output.
//
// If name is blank this method is a no-op; names longer than 128 characters
// are truncated.
func (p *Podcast) AddPerson(name, role, group, img, href string) {
	person := newPPerson(name, role, group, img, href)
	if person == nil {
		return
	}
	p.PPersons = append(p.PPersons, person)
}

// AddTranscript appends a podcast:transcript reference to the episode,
// linking a transcript or closed-caption file of the given MIME type (for
// example "text/vtt" or "application/json").
//
// language is optional and defaults to the feed language when absent. rel is
// optional; a value of "captions" marks the file as closed captions with
// timecodes.
//
// If url or mimeType is empty this method is a no-op.
func (i *Item) AddTranscript(url, mimeType, language, rel string) {
	if len(url) == 0 || len(mimeType) == 0 {
		return
	}
	i.PTranscripts = append(i.PTranscripts, &PTranscript{
		URL:      url,
		Type:     mimeType,
		Language: language,
		Rel:      rel,
	})
}

// AddChapters sets the podcast:chapters reference for the episode, replacing
// any previous value. mimeType is typically "application/json+chapters".
//
// If url or mimeType is empty this method is a no-op.
func (i *Item) AddChapters(url, mimeType string) {
	if len(url) == 0 || len(mimeType) == 0 {
		return
	}
	i.PChapters = &PChapters{
		URL:  url,
		Type: mimeType,
	}
}

// AddPerson appends a podcast:person credit to the episode. When present,
// episode-level person tags replace all channel-level people data for that
// episode, so include everyone relevant.
//
// The parameters follow the same rules as Podcast.AddPerson.
func (i *Item) AddPerson(name, role, group, img, href string) {
	person := newPPerson(name, role, group, img, href)
	if person == nil {
		return
	}
	i.PPersons = append(i.PPersons, person)
}

// AddSocialInteract appends a podcast:socialInteract reference designating
// the root post where comments and discussion for the episode take place.
//
// protocol names the interaction protocol (for example "activitypub").
// accountID identifies the account that created the post; it is recommended
// but optional.
//
// If uri or protocol is empty this method is a no-op.
func (i *Item) AddSocialInteract(uri, protocol, accountID string) {
	if len(uri) == 0 || len(protocol) == 0 {
		return
	}
	i.PSocialInteracts = append(i.PSocialInteracts, &PSocialInteract{
		URI:       uri,
		Protocol:  protocol,
		AccountID: accountID,
	})
}

func newPPerson(name, role, group, img, href string) *PPerson {
	if len(strings.TrimSpace(name)) == 0 {
		return nil
	}
	return &PPerson{
		Role:  role,
		Group: group,
		Img:   img,
		Href:  href,
		Name:  truncateRunes(name, personNameRuneLimit),
	}
}

// hasPodcastElements reports whether any Podcasting 2.0 tag is set, which
// determines whether Encode declares the podcast namespace on <rss>.
func (p *Podcast) hasPodcastElements() bool {
	if len(p.PGUID) > 0 || len(p.PMedium) > 0 || p.PLocked != nil || len(p.PPersons) > 0 {
		return true
	}
	for _, i := range p.Items {
		if i == nil {
			continue
		}
		if len(i.PTranscripts) > 0 || i.PChapters != nil ||
			len(i.PPersons) > 0 || len(i.PSocialInteracts) > 0 {
			return true
		}
	}
	return false
}
