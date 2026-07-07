package podcast_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/hbmartin/podcast-rss-generator/v2"
)

func FuzzPodcastEncode(f *testing.F) {
	invalidUTF8 := string([]byte{0xff, 0xfe, '<', '&'})
	longDescription := strings.Repeat("episode ", 1500)
	dateBytes, err := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC).GobEncode()
	if err != nil {
		f.Fatalf("encode seed date: %v", err)
	}

	f.Add(
		"Example Show",
		"https://example.com",
		"show description",
		"Episode 1",
		"episode description",
		"https://example.com/episodes/1",
		"https://example.com/episodes/1.mp3",
		int(podcast.MP3),
		int64(12345),
		dateBytes,
	)
	f.Add(
		"Symbols <>&",
		"https://example.com/feed.xml",
		`show <a href="https://example.com">site</a>`,
		"Episode <>&",
		`item <a href="https://example.com/1">site</a>`,
		"",
		"https://example.com/episodes/1.m4a",
		int(podcast.M4A),
		int64(-1),
		[]byte{},
	)
	f.Add(
		invalidUTF8,
		"",
		longDescription,
		invalidUTF8,
		longDescription,
		"https://example.com/episodes/invalid",
		"",
		int(podcast.EnclosureUnknown),
		int64(0),
		[]byte{0xff, 0x00, 0x01},
	)

	f.Fuzz(func(
		_ *testing.T,
		title string,
		link string,
		description string,
		itemTitle string,
		itemDescription string,
		itemLink string,
		enclosureURL string,
		enclosureType int,
		enclosureLength int64,
		dateBytes []byte,
	) {
		p := podcast.New(title, link, description, time.Time{}, time.Time{})
		p.AddAtomLink(link)
		p.AddCategory("Arts", []string{itemTitle, itemDescription})
		p.AddSummary(description)
		p.SetPodcastGUID(podcast.NewFeedGUID(link))
		p.SetLocked(true, itemTitle)
		p.AddPerson(itemTitle, description, "", enclosureURL, link)

		if date, ok := fuzzDate(dateBytes); ok {
			p.AddPubDate(date)
			p.AddLastBuildDate(date)
		}

		i := podcast.Item{
			Title:       itemTitle,
			Description: podcast.Description(itemDescription),
			Link:        itemLink,
		}
		i.AddImage(enclosureURL)
		i.AddSummary(itemDescription)
		i.AddTranscript(enclosureURL, "text/vtt", itemTitle, "captions")
		i.AddChapters(itemLink, "application/json+chapters")
		i.AddPerson(itemDescription, itemTitle, "", "", itemLink)
		i.AddSocialInteract(itemLink, "activitypub", itemTitle)
		if date, ok := fuzzDate(dateBytes); ok {
			i.AddPubDate(date)
		}
		if len(enclosureURL) != 0 {
			i.AddEnclosure(enclosureURL, podcast.EnclosureType(enclosureType), enclosureLength)
		}

		if _, err := p.AddItem(i); err != nil {
			return
		}

		var buf bytes.Buffer
		if err := p.Encode(&buf); err != nil {
			return
		}
	})
}

func fuzzDate(data []byte) (time.Time, bool) {
	var t time.Time
	if err := t.GobDecode(data); err != nil {
		return time.Time{}, false
	}
	return t, true
}
