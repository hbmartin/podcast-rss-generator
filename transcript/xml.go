package transcript

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// pstNamespace is the Podlove simple-transcripts namespace that gates XML
// transcript parsing.
const pstNamespace = "http://podlove.org/simple-transcripts"

// errMissingItemAttr reports a transcript <item> missing a timing attribute.
var errMissingItemAttr = errors.New("transcript item missing timing attribute")

// xmlNode is a minimal XML element tree keyed by local element names.
type xmlNode struct {
	Local    string
	Attrs    map[string]string
	Children []*xmlNode
	Text     string
}

// fullText returns the concatenated text content of the node and descendants,
// mirroring BeautifulSoup's Tag.text for the leaf items used here.
func (n *xmlNode) fullText() string {
	var b strings.Builder
	b.WriteString(n.Text)
	for _, child := range n.Children {
		b.WriteString(child.fullText())
	}
	return b.String()
}

// ParsePodloveXML parses a Podlove simple-transcripts XML document into a
// Transcript. It returns an *InvalidXMLError when the document is not a Podlove
// transcript, or a *NoTranscriptFoundError when it contains no transcript.
func ParsePodloveXML(xmlString string) (*Transcript, error) {
	if !strings.Contains(xmlString, pstNamespace) {
		return nil, newInvalidXMLError()
	}
	root, err := parseXMLTree([]byte(xmlString))
	if err != nil {
		return nil, newInvalidXMLError()
	}
	segments, err := xmlToSegments(root)
	if err != nil {
		return nil, err
	}
	return newTranscript(segments), nil
}

func xmlToSegments(root *xmlNode) ([]Segment, error) {
	transcripts := findFirst(root, "transcripts")
	if transcripts == nil {
		return nil, newNoTranscriptFoundError()
	}

	segments := []*Segment{{}}
	for _, speech := range transcripts.Children {
		if speech.Local != "speech" {
			continue
		}
		last := segments[len(segments)-1]
		if last.Body != nil {
			last = &Segment{}
			segments = append(segments, last)
		}
		for _, item := range speech.Children {
			if item.Local != "item" {
				continue
			}
			if err := appendItem(last, item); err != nil {
				return nil, err
			}
		}
	}

	populated := make([]Segment, 0, len(segments))
	for _, segment := range segments {
		if !segment.isEmpty() {
			populated = append(populated, *segment)
		}
	}
	if len(populated) == 0 {
		return nil, newNoTranscriptFoundError()
	}
	return populated, nil
}

// appendItem folds one <item> into the current segment.
func appendItem(segment *Segment, item *xmlNode) error {
	text := item.fullText()
	if segment.Body == nil {
		segment.Body = strPtr(text)
	} else {
		*segment.Body += " " + text
	}

	start, ok := item.Attrs["start"]
	if !ok {
		return fmt.Errorf("%w: start", errMissingItemAttr)
	}
	end, ok := item.Attrs["end"]
	if !ok {
		return fmt.Errorf("%w: end", errMissingItemAttr)
	}
	if segment.StartTime == nil {
		startSecs, err := mtsToSecsFloat(start)
		if err != nil {
			return err
		}
		segment.StartTime = floatPtr(startSecs)
	}
	endSecs, err := mtsToSecsFloat(end)
	if err != nil {
		return err
	}
	segment.EndTime = floatPtr(endSecs)
	return nil
}

// findFirst returns the first element (depth-first, self included) with the
// given local name, or nil.
func findFirst(node *xmlNode, local string) *xmlNode {
	if node.Local == local {
		return node
	}
	for _, child := range node.Children {
		if found := findFirst(child, local); found != nil {
			return found
		}
	}
	return nil
}

// parseXMLTree builds an xmlNode tree from XML bytes, keying elements and
// attributes by their local names.
func parseXMLTree(data []byte) (*xmlNode, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var root *xmlNode
	var stack []*xmlNode
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch tok := token.(type) {
		case xml.StartElement:
			node := &xmlNode{Local: tok.Name.Local, Attrs: make(map[string]string, len(tok.Attr))}
			for _, a := range tok.Attr {
				if a.Name.Local == "xmlns" || a.Name.Space == "xmlns" {
					continue
				}
				node.Attrs[a.Name.Local] = a.Value
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, node)
			} else if root == nil {
				root = node
			}
			stack = append(stack, node)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			if len(stack) > 0 {
				stack[len(stack)-1].Text += string(tok)
			}
		}
	}
	if root == nil {
		return nil, io.ErrUnexpectedEOF
	}
	return root, nil
}

// XMLFileToJSONFile converts a Podlove XML file to PodcastIndex transcript JSON,
// merging optional metadata. Nothing is written when parsing fails.
func XMLFileToJSONFile(xmlFile, jsonFile string, metadata map[string]string) error {
	xmlString, err := ReadTextRobust(xmlFile)
	if err != nil {
		return err
	}
	transcript, err := ParsePodloveXML(xmlString)
	if err != nil {
		var invalidXML *InvalidXMLError
		if errors.As(err, &invalidXML) {
			return fmt.Errorf("%w: %s", err, xmlFile)
		}
		return err
	}
	if len(metadata) > 0 {
		transcript.Metadata = metadata
	}
	return writeTranscriptJSON(jsonFile, transcript)
}
