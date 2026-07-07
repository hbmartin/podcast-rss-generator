package podcast_test

import (
	"encoding/xml"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2"
)

func TestDescriptionMarshalXML(t *testing.T) {
	t.Parallel()

	type document struct {
		XMLName     xml.Name            `xml:"doc"`
		Description podcast.Description `xml:"description"`
	}

	tests := []struct {
		name string
		in   podcast.Description
		want string
	}{
		{
			name: "plain text stays text",
			in:   "plain description",
			want: `<doc><description>plain description</description></doc>`,
		},
		{
			name: "rich text uses cdata",
			in:   `read <a href="https://example.com">more</a>`,
			want: `<doc><description><![CDATA[read <a href="https://example.com">more</a>]]></description></doc>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := xml.Marshal(document{Description: tt.in})
			if err != nil {
				t.Fatalf("xml.Marshal(Description) error = %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("xml.Marshal(Description) = %q, want %q", got, tt.want)
			}
		})
	}
}
