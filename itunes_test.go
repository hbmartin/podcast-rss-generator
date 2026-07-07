package podcast_test

import (
	"encoding/xml"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
)

func TestITunesTypesMarshalXML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   any
		want string
	}{
		{
			name: "category with nested category",
			in: &podcast.ICategory{
				Text: "Technology",
				ICategories: []*podcast.ICategory{
					{Text: "Podcasting"},
				},
			},
			want: `<itunes:category text="Technology"><itunes:category text="Podcasting"></itunes:category></itunes:category>`,
		},
		{
			name: "image href",
			in:   &podcast.IImage{HREF: "https://example.com/itunes.jpg"},
			want: `<itunes:image href="https://example.com/itunes.jpg"></itunes:image>`,
		},
		{
			name: "summary cdata",
			in:   &podcast.ISummary{Text: `summary <a href="https://example.com">link</a>`},
			want: `<itunes:summary><![CDATA[summary <a href="https://example.com">link</a>]]></itunes:summary>`,
		},
		{
			name: "podcast type",
			in:   &podcast.IType{Text: "serial"},
			want: `<itunes:type>serial</itunes:type>`,
		},
		{
			name: "episode type",
			in:   &podcast.IEpisodeType{Text: "bonus"},
			want: `<itunes:episodeType>bonus</itunes:episodeType>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := xml.Marshal(tt.in)
			if err != nil {
				t.Fatalf("xml.Marshal(%T) error = %v", tt.in, err)
			}
			if string(got) != tt.want {
				t.Fatalf("xml.Marshal(%T) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
