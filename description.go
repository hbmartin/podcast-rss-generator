package podcast

import "encoding/xml"

// AtomLink represents the Atom reference link.
type Description struct {
	XMLName xml.Name `xml:"description"`
	Text    string   `xml:",cdata"`
}
