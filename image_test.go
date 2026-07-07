package podcast_test

import (
	"encoding/xml"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
)

func TestImageMarshalXML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		image podcast.Image
		want  string
	}{
		{
			name: "required fields only omit zero values",
			image: podcast.Image{
				URL:   "https://example.com/art.jpg",
				Title: "Example Show",
				Link:  "https://example.com",
			},
			want: `<image><url>https://example.com/art.jpg</url><title>Example Show</title><link>https://example.com</link></image>`,
		},
		{
			name: "optional dimensions and description",
			image: podcast.Image{
				URL:         "https://example.com/art.jpg",
				Title:       "Example Show",
				Link:        "https://example.com",
				Description: "cover art",
				Width:       1400,
				Height:      1400,
			},
			want: `<image><url>https://example.com/art.jpg</url><title>Example Show</title><link>https://example.com</link><description>cover art</description><width>1400</width><height>1400</height></image>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := xml.Marshal(tt.image)
			if err != nil {
				t.Fatalf("xml.Marshal(Image) error = %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("xml.Marshal(Image) = %q, want %q", got, tt.want)
			}
		})
	}
}
