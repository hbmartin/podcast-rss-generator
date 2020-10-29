package podcast

import (
	"encoding/xml"
	"strings"
)

// Description is a 4000 character rich-text field for the channel and
// podcast description tags.
//
// This is rendered as CDATA which allows for HTML tags such as `<a href="">`.
type Description string

// MarshalXML renders the description text as CDATA so rich text such as
// `<a href="">` stays valid XML.
func (d Description) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !strings.ContainsAny(string(d), "<>&") {
		return e.EncodeElement(string(d), start)
	}

	v := struct {
		Text string `xml:",cdata"`
	}{
		Text: string(d),
	}
	return e.EncodeElement(v, start)
}
