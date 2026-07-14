package chapters

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

// attr is a single XML attribute (local name, or "{ns}local" when namespaced).
type attr struct {
	Key   string
	Value string
}

// Element is a minimal, read-only XML element tree mirroring the subset of
// xml.etree.ElementTree that the chapter extractors rely on. Tag uses Clark
// notation ("{namespace}local", or just "local" when unqualified), matching
// ElementTree, so callers can navigate namespaced feeds the same way.
type Element struct {
	Tag      string
	Attrs    []attr
	Children []*Element
	Text     string
}

// clark returns the Clark-notation form of an xml.Name.
func clark(name xml.Name) string {
	if name.Space == "" {
		return name.Local
	}
	return "{" + name.Space + "}" + name.Local
}

// Get returns the value of the named attribute and whether it was present.
func (e *Element) Get(key string) (string, bool) {
	for _, a := range e.Attrs {
		if a.Key == key {
			return a.Value, true
		}
	}
	return "", false
}

// Find returns the first child element matching the given tag, which may be a
// Clark-notation tag and may carry a leading "./" path prefix (the only XPath
// forms the extractors use). It returns nil when no child matches.
func (e *Element) Find(tag string) *Element {
	tag = strings.TrimPrefix(tag, "./")
	for _, child := range e.Children {
		if child.Tag == tag {
			return child
		}
	}
	return nil
}

// ParseXML parses an XML document into an Element tree. It returns an error for
// malformed XML.
func ParseXML(data []byte) (*Element, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	var root *Element
	stack := []*Element{}
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
			element := &Element{Tag: clark(tok.Name)}
			for _, a := range tok.Attr {
				if a.Name.Space == "xmlns" || a.Name.Local == "xmlns" {
					continue
				}
				element.Attrs = append(element.Attrs, attr{Key: clark(a.Name), Value: a.Value})
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, element)
			} else if root == nil {
				root = element
			}
			stack = append(stack, element)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			// Mirror ElementTree's .text: character data before the first
			// child element of the current element.
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				if len(current.Children) == 0 {
					current.Text += string(tok)
				}
			}
		}
	}
	if root == nil {
		return nil, io.ErrUnexpectedEOF
	}
	return root, nil
}
