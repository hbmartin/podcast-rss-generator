package chapters

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// DefaultTimeout is the default HTTP timeout for the feed- and chapters-fetching
// helpers.
const DefaultTimeout = 30 * time.Second

// archiveFilePerm is used for the best-effort PCI chapters cache file.
const archiveFilePerm = 0o600

// minChapters is the fewest chapters a description must yield to be treated as a
// real chapter list.
const minChapters = 2

// HTTPClient is the subset of *http.Client used to fetch feeds and chapter
// documents. *http.Client satisfies it; tests inject fakes.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// fetchConfig holds resolved options for the fetching helpers.
type fetchConfig struct {
	headers     map[string]string
	timeout     time.Duration
	archivePath string
}

// FetchOption configures the HTTP fetching helpers.
type FetchOption func(*fetchConfig)

// WithHeaders sets request headers for a fetch.
func WithHeaders(headers map[string]string) FetchOption {
	return func(c *fetchConfig) { c.headers = headers }
}

// WithTimeout overrides the default request timeout.
func WithTimeout(timeout time.Duration) FetchOption {
	return func(c *fetchConfig) { c.timeout = timeout }
}

// WithArchivePath caches the fetched PodcastIndex chapters JSON at path so
// repeated calls read it instead of refetching. It only affects
// GetAndExtractPCIChapters.
func WithArchivePath(path string) FetchOption {
	return func(c *fetchConfig) { c.archivePath = path }
}

// ExtractDescriptionChapters extracts chapters from an episode description
// (plain text or HTML). It returns nil when fewer than two chapters are found,
// since a lone timestamp is more likely an incidental mention than a chapter
// list. When stripHTML is true, markup is removed from chapter titles.
func ExtractDescriptionChapters(description string, stripHTML bool) []Chapter {
	descChapters := findAllPairs(descriptionChapterRe, description)
	if chapters := descMatchesToChapters(descChapters, stripHTML); chapters != nil {
		return chapters
	}
	if len(descChapters) == 1 {
		retry := findAllPairs(retryDescriptionChapterRe, descChapters[0][0]+" "+descChapters[0][1])
		return descMatchesToChapters(retry, stripHTML)
	}
	return nil
}

func descMatchesToChapters(matches [][2]string, stripHTML bool) []Chapter {
	if len(matches) < minChapters {
		return nil
	}
	chapters := make([]Chapter, 0, len(matches))
	for _, match := range matches {
		chapter, err := extractDescTsAndTitle(match, stripHTML)
		if err != nil {
			continue
		}
		chapters = append(chapters, chapter)
	}
	if len(chapters) >= minChapters {
		return chapters
	}
	return nil
}

func extractDescTsAndTitle(tsTitle [2]string, stripHTML bool) (Chapter, error) {
	title := strings.TrimSpace(tsTitle[1])
	if strings.Contains(title, "<a") && !strings.Contains(title, "</a>") {
		title += "</a>"
	}
	url := urlRe.FindString(title)
	if stripHTML {
		title = StripHTML(title)
	}
	seconds, err := TsToSecs(tsTitle[0])
	if err != nil {
		return Chapter{}, err
	}
	return Chapter{Start: seconds, Title: title, URL: url}, nil
}

// ExtractPCIChapters extracts chapters from a parsed PodcastIndex chapters JSON
// document. It returns nil when the document is malformed.
func ExtractPCIChapters(chaptersJSON map[string]any) []Chapter {
	rawChapters, ok := chaptersJSON["chapters"].([]any)
	if !ok {
		return nil
	}
	chapters := make([]Chapter, 0, len(rawChapters))
	for _, raw := range rawChapters {
		entry, ok := raw.(map[string]any)
		if !ok {
			return nil
		}
		start, ok := coerceInt(entry["startTime"])
		if !ok {
			return nil
		}
		title, ok := entry["title"].(string)
		if !ok {
			return nil
		}
		chapters = append(chapters, Chapter{
			Start: start,
			Title: title,
			URL:   optionalString(entry, "url"),
			Image: optionalString(entry, "img"),
		})
	}
	return chapters
}

// GetAndExtractPCIChapters fetches a PodcastIndex chapters JSON document and
// extracts its chapters, returning nil on any error. With WithArchivePath, the
// raw JSON is read from (or written to) that path so repeated calls do not
// refetch the document.
func GetAndExtractPCIChapters(ctx context.Context, client HTTPClient, url string, opts ...FetchOption) []Chapter {
	cfg := fetchConfig{timeout: DefaultTimeout}
	for _, opt := range opts {
		opt(&cfg)
	}

	var parsed any
	if cfg.archivePath != "" && fileExists(cfg.archivePath) {
		data, err := os.ReadFile(cfg.archivePath)
		if err != nil {
			return nil
		}
		if !utf8.Valid(data) {
			return nil
		}
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil
		}
	} else {
		body, ok := fetch(ctx, client, url, cfg)
		if !ok {
			return nil
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil
		}
		if cfg.archivePath != "" {
			writeArchive(cfg.archivePath, parsed)
		}
	}

	doc, ok := parsed.(map[string]any)
	if !ok {
		return nil
	}
	return ExtractPCIChapters(doc)
}

// writeArchive best-effort persists the parsed chapters JSON; write failures
// are intentionally ignored so extraction still returns.
func writeArchive(path string, parsed any) {
	serialized, err := json.Marshal(parsed)
	if err != nil {
		return
	}
	if writeErr := os.WriteFile(path, serialized, archiveFilePerm); writeErr != nil {
		_ = writeErr
	}
}

// ExtractPSCChaptersFromFile extracts PSC chapters for the episode with guid
// from a feed file, returning nil when unavailable.
func ExtractPSCChaptersFromFile(feedFile, guid string) []Chapter {
	root := readAndParseFeedFile(feedFile)
	if root == nil {
		return nil
	}
	return extractPSCChaptersForGUID(root, guid)
}

// ExtractPSCChaptersFromURL fetches a podcast feed and extracts PSC chapters
// for guid, returning nil when unavailable.
func ExtractPSCChaptersFromURL(ctx context.Context, client HTTPClient, feedURL, guid string, opts ...FetchOption) []Chapter {
	cfg := fetchConfig{timeout: DefaultTimeout}
	for _, opt := range opts {
		opt(&cfg)
	}
	body, ok := fetch(ctx, client, feedURL, cfg)
	if !ok {
		return nil
	}
	root := parseFeed(body)
	if root == nil {
		return nil
	}
	return extractPSCChaptersForGUID(root, guid)
}

func extractPSCChaptersForGUID(root *Element, guid string) []Chapter {
	for _, item := range iterFeedItems(root) {
		foundGUID := item.Find("guid")
		if foundGUID != nil && foundGUID.Text == guid {
			if pscChapters := item.Find("./{" + pscNamespace + "}chapters"); pscChapters != nil {
				return ExtractPSCChapters(pscChapters)
			}
			return nil
		}
	}
	return nil
}

// ExtractAllPSCChaptersFromFile extracts PSC chapters for every episode in a
// feed file, keyed by GUID. It returns nil when the feed cannot be read or
// parsed.
func ExtractAllPSCChaptersFromFile(feedFile string) map[string][]Chapter {
	root := readAndParseFeedFile(feedFile)
	if root == nil {
		return nil
	}
	allChapters := map[string][]Chapter{}
	for _, item := range iterFeedItems(root) {
		foundGUID := item.Find("guid")
		if foundGUID == nil {
			continue
		}
		if pscChapters := item.Find("./{" + pscNamespace + "}chapters"); pscChapters != nil {
			if chapters := ExtractPSCChapters(pscChapters); chapters != nil {
				allChapters[foundGUID.Text] = chapters
			}
		}
	}
	return allChapters
}

// ExtractPSCChapters extracts PSC chapters from a <psc:chapters> element. It
// returns nil when any chapter lacks a start or title, or has an invalid start.
func ExtractPSCChapters(pscChapters *Element) []Chapter {
	chapters := make([]Chapter, 0, len(pscChapters.Children))
	for _, child := range pscChapters.Children {
		start, ok := child.Get("start")
		if !ok {
			return nil
		}
		title, ok := child.Get("title")
		if !ok {
			return nil
		}
		seconds, err := TsToSecs(start)
		if err != nil {
			return nil
		}
		href, _ := child.Get("href")
		image, _ := child.Get("image")
		chapters = append(chapters, Chapter{
			Start: seconds,
			Title: title,
			URL:   href,
			Image: image,
		})
	}
	return chapters
}

// FindPCIChaptersURL returns the <podcast:chapters> URL declared for the
// episode with guid, ready to pass to GetAndExtractPCIChapters, or "" when not
// declared.
func FindPCIChaptersURL(feedFile, guid string) string {
	root := readAndParseFeedFile(feedFile)
	if root == nil {
		return ""
	}
	for _, item := range iterFeedItems(root) {
		foundGUID := item.Find("guid")
		if foundGUID != nil && foundGUID.Text == guid {
			if pciChapters := item.Find("./{" + pciNamespace + "}chapters"); pciChapters != nil {
				url, _ := pciChapters.Get("url")
				return url
			}
			return ""
		}
	}
	return ""
}

// iterFeedItems returns the <item> elements of the feed's channel.
func iterFeedItems(root *Element) []*Element {
	channel := root.Find("./channel")
	if channel == nil {
		return nil
	}
	var items []*Element
	for _, element := range channel.Children {
		if element.Tag == "item" {
			items = append(items, element)
		}
	}
	return items
}

// readAndParseFeedFile reads a feed file and parses it, returning nil when the
// file is missing, is not valid UTF-8, is unparseable, or has no channel.
func readAndParseFeedFile(feedFile string) *Element {
	if !fileExists(feedFile) {
		return nil
	}
	data, err := os.ReadFile(feedFile) //nolint:gosec // reading a caller-supplied feed path is this function's purpose
	if err != nil {
		return nil
	}
	if !utf8.Valid(data) {
		return nil
	}
	return parseFeed(data)
}

// parseFeed parses feed XML and returns the root, or nil when the XML is
// malformed or has no channel element.
func parseFeed(data []byte) *Element {
	root, err := ParseXML(data)
	if err != nil {
		return nil
	}
	if root.Find("./channel") == nil {
		return nil
	}
	return root
}

// fetch performs an HTTP GET and returns the body and whether it succeeded
// (a network error or a status >= 400 yields ok=false).
func fetch(ctx context.Context, client HTTPClient, url string, cfg fetchConfig) ([]byte, bool) {
	reqCtx := ctx
	if cfg.timeout > 0 {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(ctx, cfg.timeout)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false
	}
	for key, value := range cfg.headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false
	}
	//nolint:errcheck // the response body close error is not actionable here
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false
	}
	return body, true
}

// coerceInt converts a JSON value to an int the way Python's int() does for the
// values seen here: numbers truncate toward zero and integer strings parse.
func coerceInt(value any) (int, bool) {
	switch n := value.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return int(i), true
		}
		if f, err := n.Float64(); err == nil {
			return int(f), true
		}
		return 0, false
	case string:
		i, err := strconv.Atoi(strings.TrimSpace(n))
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}

// optionalString returns entry[key] when it is a string, else "".
func optionalString(entry map[string]any, key string) string {
	if value, ok := entry[key].(string); ok {
		return value
	}
	return ""
}

// fileExists reports whether path exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
