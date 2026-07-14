package transcript

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// HTML transcript detection regexes, mirroring the Python source.
var (
	htmlPTsRe = regexp.MustCompile(`^<p>\s*(\d+:)?\d+:\d+\s*</p>`)
	htmlTsRe  = regexp.MustCompile(`^\s*(\d+:)?\d+:\d+\s*$`)
)

// secondsPerHTMLUnit is the base (60) for HTML colon-separated timestamps.
const secondsPerHTMLUnit = 60

// htmlTsToSecs converts a colon-separated timestamp ([[HH:]MM:]SS) to seconds.
func htmlTsToSecs(timeString string) (int, error) {
	parts := strings.Split(timeString, ":")
	total := 0
	for i, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return 0, err
		}
		total += intPow(secondsPerHTMLUnit, len(parts)-1-i) * value
	}
	return total, nil
}

// intPow returns base**exp for small non-negative exponents.
func intPow(base, exp int) int {
	result := 1
	for range exp {
		result *= base
	}
	return result
}

// ParseHTML parses an HTML transcript into a Transcript. It returns an
// *InvalidHTMLError when the format gate fails, or a *NoTranscriptFoundError
// when no transcript content is present.
func ParseHTML(htmlString string) (*Transcript, error) {
	if !strings.Contains(htmlString, "<cite>") &&
		!strings.Contains(htmlString, "<time>") &&
		!htmlPTsRe.MatchString(htmlString) {
		return nil, newInvalidHTMLError()
	}
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, newInvalidHTMLError()
	}
	segments, err := htmlToSegments(doc)
	if err != nil {
		return nil, err
	}
	return newTranscript(segments), nil
}

func htmlToSegments(doc *html.Node) ([]Segment, error) {
	body := findElement(doc, "body")
	segments := []*Segment{{}}
	iterate := body
	if iterate == nil {
		iterate = doc
	}
	for child := iterate.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}
		if err := applyHTMLNode(&segments, child); err != nil {
			return nil, err
		}
	}
	if segments[0].isEmpty() {
		return nil, newNoTranscriptFoundError()
	}
	result := make([]Segment, len(segments))
	for i, segment := range segments {
		result[i] = *segment
	}
	return result, nil
}

// applyHTMLNode folds one <cite>, <time>, or <p> element into the segment list.
func applyHTMLNode(segments *[]*Segment, node *html.Node) error {
	switch node.Data {
	case "cite":
		speaker := strings.TrimSpace(strings.ReplaceAll(nodeText(node), ":", ""))
		addSpeaker(segments, speaker)
	case "time":
		seconds, err := htmlTsToSecs(strings.TrimSpace(nodeText(node)))
		if err != nil {
			return err
		}
		addStartTime(segments, seconds)
	case "p":
		stripped := strings.TrimSpace(nodeText(node))
		if stripped == "" {
			return nil
		}
		if htmlTsRe.MatchString(stripped) {
			seconds, err := htmlTsToSecs(stripped)
			if err != nil {
				return err
			}
			addStartTime(segments, seconds)
		} else {
			(*segments)[len(*segments)-1].Body = strPtr(stripped)
		}
	}
	return nil
}

func addSpeaker(segments *[]*Segment, speaker string) {
	last := (*segments)[len(*segments)-1]
	if last.Speaker == nil {
		last.Speaker = strPtr(speaker)
	} else {
		*segments = append(*segments, &Segment{Speaker: strPtr(speaker)})
	}
}

func addStartTime(segments *[]*Segment, startTime int) {
	last := (*segments)[len(*segments)-1]
	if last.StartTime == nil {
		last.StartTime = floatPtr(float64(startTime))
	} else {
		*segments = append(*segments, &Segment{StartTime: floatPtr(float64(startTime))})
	}
}

// findElement returns the first element node with the given tag name, or nil.
func findElement(node *html.Node, name string) *html.Node {
	if node.Type == html.ElementNode && node.Data == name {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findElement(child, name); found != nil {
			return found
		}
	}
	return nil
}

// nodeText returns the concatenated text of a node and its descendants.
func nodeText(node *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return b.String()
}

// HTMLFileToJSONFile converts an HTML file to PodcastIndex transcript JSON,
// merging optional metadata. Nothing is written when parsing fails.
func HTMLFileToJSONFile(htmlFile, jsonFile string, metadata map[string]string) error {
	htmlString, err := ReadTextRobust(htmlFile)
	if err != nil {
		return err
	}
	transcript, err := ParseHTML(htmlString)
	if err != nil {
		return fmt.Errorf("%w: %s", err, htmlFile)
	}
	if len(metadata) > 0 {
		transcript.Metadata = metadata
	}
	return writeTranscriptJSON(jsonFile, transcript)
}
