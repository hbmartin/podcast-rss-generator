package podcast_test

import (
	"encoding/xml"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
)

func TestTextInputMarshalXML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   podcast.TextInput
		want string
	}{
		{
			name: "rss text input",
			in: podcast.TextInput{
				Title:       "Search",
				Description: "Search the archive",
				Name:        "q",
				Link:        "https://example.com/search",
			},
			want: `<textInput><title>Search</title><description>Search the archive</description><name>q</name><link>https://example.com/search</link></textInput>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := xml.Marshal(tt.in)
			if err != nil {
				t.Fatalf("xml.Marshal(TextInput) error = %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("xml.Marshal(TextInput) = %q, want %q", got, tt.want)
			}
		})
	}
}
