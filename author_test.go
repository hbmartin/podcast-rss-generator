package podcast_test

import (
	"encoding/xml"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
)

func TestAuthorMarshalXML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		author podcast.Author
		want   string
	}{
		{
			name: "itunes owner fields",
			author: podcast.Author{
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
			want: `<Author><itunes:name>Jane Doe</itunes:name><itunes:email>jane@example.com</itunes:email></Author>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := xml.Marshal(tt.author)
			if err != nil {
				t.Fatalf("xml.Marshal(Author) error = %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("xml.Marshal(Author) = %q, want %q", got, tt.want)
			}
		})
	}
}
