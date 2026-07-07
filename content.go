package podcast

import "encoding/xml"

// contentNamespace is the XML namespace for the RSS content module, declared on
// the <rss> element whenever an item carries a <content:encoded> tag.
const contentNamespace = "http://purl.org/rss/1.0/modules/content/"

// EncodedContent renders a <content:encoded> element from the RSS content
// module (http://purl.org/rss/1.0/modules/content/).
//
// It carries the full, rich-HTML body of an episode. Where <description> is
// rendered as plain text by many aggregators, content:encoded is defined to
// carry markup, so clients such as Apple Podcasts display the formatted HTML –
// links, lists and paragraphs – as authored.
//
// This is rendered as CDATA which allows for HTML tags such as `<a href="">`.
type EncodedContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Text    string   `xml:",cdata"`
}
