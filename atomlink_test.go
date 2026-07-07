package podcast_test

import (
	"encoding/xml"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
)

func TestAtomLinkMarshalXML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		link podcast.AtomLink
		want string
	}{
		{
			name: "self rss link",
			link: podcast.AtomLink{
				HREF: "https://example.com/feed.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			want: `<atom:link href="https://example.com/feed.xml" rel="self" type="application/rss+xml"></atom:link>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := xml.Marshal(tt.link)
			if err != nil {
				t.Fatalf("xml.Marshal(AtomLink) error = %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("xml.Marshal(AtomLink) = %q, want %q", got, tt.want)
			}
		})
	}
}
