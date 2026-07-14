package chapters

import (
	"bytes"
	"encoding/json"
	"strings"
)

// Serialization format versions.
const (
	pciVersion = "1.2.0"
	pscVersion = "1.2"
)

// PCIChapter is a single chapter entry of a PodcastIndex chapters document.
type PCIChapter struct {
	StartTime int    `json:"startTime"`
	Title     string `json:"title"`
	URL       string `json:"url,omitempty"`
	Img       string `json:"img,omitempty"`
}

// PCIDocument is a PodcastIndex chapters document.
type PCIDocument struct {
	Version  string       `json:"version"`
	Chapters []PCIChapter `json:"chapters"`
}

// ToPCIDict builds a PodcastIndex chapters document from chapters.
func ToPCIDict(chapters []Chapter) PCIDocument {
	pciChapters := make([]PCIChapter, 0, len(chapters))
	for _, chapter := range chapters {
		pciChapters = append(pciChapters, PCIChapter{
			StartTime: chapter.Start,
			Title:     chapter.Title,
			URL:       chapter.URL,
			Img:       chapter.Image,
		})
	}
	return PCIDocument{Version: pciVersion, Chapters: pciChapters}
}

// ToPCIJSON serializes chapters to a PodcastIndex chapters JSON string.
// When indent is greater than zero the output is pretty-printed with that many
// spaces per level; otherwise it is compact.
func ToPCIJSON(chapters []Chapter, indent int) (string, error) {
	return marshalJSON(ToPCIDict(chapters), indent)
}

// marshalJSON encodes v as JSON without HTML-escaping (matching Python's
// json.dumps defaults), optionally indented.
func marshalJSON(v any, indent int) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if indent > 0 {
		encoder.SetIndent("", strings.Repeat(" ", indent))
	}
	if err := encoder.Encode(v); err != nil {
		return "", err
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// ToPSCElement builds a <psc:chapters> element containing chapters.
func ToPSCElement(chapters []Chapter) (*Element, error) {
	root := &Element{
		Tag:   "{" + pscNamespace + "}chapters",
		Attrs: []attr{{Key: "version", Value: pscVersion}},
	}
	for _, chapter := range chapters {
		start, err := SecsToTs(chapter.Start)
		if err != nil {
			return nil, err
		}
		child := &Element{
			Tag:   "{" + pscNamespace + "}chapter",
			Attrs: []attr{{Key: "start", Value: start}, {Key: "title", Value: chapter.Title}},
		}
		if chapter.URL != "" {
			child.Attrs = append(child.Attrs, attr{Key: "href", Value: chapter.URL})
		}
		if chapter.Image != "" {
			child.Attrs = append(child.Attrs, attr{Key: "image", Value: chapter.Image})
		}
		root.Children = append(root.Children, child)
	}
	return root, nil
}

// ToPSCXML serializes chapters to a Podlove Simple Chapters XML string.
func ToPSCXML(chapters []Chapter) (string, error) {
	root, err := ToPSCElement(chapters)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	writePSCElement(&b, root, true)
	return b.String(), nil
}

// writePSCElement serializes a PSC element using the "psc" namespace prefix,
// matching xml.etree.ElementTree.tostring output.
func writePSCElement(b *strings.Builder, el *Element, isRoot bool) {
	local := strings.TrimPrefix(el.Tag, "{"+pscNamespace+"}")
	b.WriteString("<psc:")
	b.WriteString(local)
	if isRoot {
		b.WriteString(` xmlns:psc="`)
		b.WriteString(pscNamespace)
		b.WriteString(`"`)
	}
	for _, a := range el.Attrs {
		b.WriteString(" ")
		b.WriteString(a.Key)
		b.WriteString(`="`)
		b.WriteString(escapeAttr(a.Value))
		b.WriteString(`"`)
	}
	if len(el.Children) == 0 {
		b.WriteString(" />")
		return
	}
	b.WriteString(">")
	for _, child := range el.Children {
		writePSCElement(b, child, false)
	}
	b.WriteString("</psc:")
	b.WriteString(local)
	b.WriteString(">")
}

// escapeAttr escapes an XML attribute value the way ElementTree does.
func escapeAttr(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"\n", "&#10;",
		"\t", "&#09;",
		"\r", "&#13;",
	)
	return replacer.Replace(value)
}

// ToDescription serializes chapters to description-embed text, one
// chapter per line. A chapter's URL is appended only when it is not already
// contained in the title.
func ToDescription(chapters []Chapter) (string, error) {
	lines := make([]string, 0, len(chapters))
	for _, chapter := range chapters {
		start, err := SecsToTs(chapter.Start)
		if err != nil {
			return "", err
		}
		line := start + " " + chapter.Title
		if chapter.URL != "" && !strings.Contains(chapter.Title, chapter.URL) {
			line += " " + chapter.URL
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
}
