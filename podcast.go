package podcast

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	pVersion                    = "2.0.0"
	iTunesSubtitleRuneLimit     = 64
	iTunesSummaryRuneLimit      = 4000
	itemDescriptionRuneLimit    = 10000
	podcastDescriptionByteLimit = 4000
)

// Sentinel errors returned by feed and item validation.
var (
	ErrPodcastRequired          = errors.New("podcast is required")
	ErrWriterRequired           = errors.New("writer is required")
	ErrTitleDescriptionRequired = errors.New("title and description are required")
	ErrEnclosureURLRequired     = errors.New("enclosure url is required")
	ErrEnclosureTypeRequired    = errors.New("enclosure type is required")
	ErrLinkRequired             = errors.New("link is required when not using enclosure")

	// Channel-level validation errors returned by Podcast.Validate.
	ErrChannelTitleRequired       = errors.New("channel title is required")
	ErrChannelDescriptionRequired = errors.New("channel description is required")
	ErrChannelLinkRequired        = errors.New("channel link is required")
)

// ItemValidationError describes why an item could not be added to a podcast.
type ItemValidationError struct {
	Title string
	Err   error
}

// Error returns a human-readable item validation message.
func (e *ItemValidationError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return "item validation failed"
	}
	if len(e.Title) == 0 {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Title, e.Err)
}

// Unwrap returns the underlying validation error.
func (e *ItemValidationError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Podcast represents a podcast.
type Podcast struct {
	XMLName xml.Name `xml:"channel"`

	// Title is the show title.
	//
	// This is a required tag.
	//
	// It’s important to have a clear, concise name for your podcast. Make your
	// title specific. A show titled Our Community Bulletin is too vague to
	// attract many subscribers, no matter how compelling the content.
	//
	// Pay close attention to the title as Apple Podcasts uses this field fo
	// search.
	//
	// If you include a long list of keywords in an attempt to game podcast
	// search, your show may be removed from the Apple directory.
	Title string `xml:"title"`

	// Link is the associated with a podcast.
	//
	// Do not specify HTML here.  Use RAW https:// urls.
	Link string `xml:"link"`

	// Description is text containing one or more sentences describing
	// your podcast to potential listeners.
	//
	// This is a required tag.
	//
	// Limit: 4000 bytes
	//
	// Note that this field is a CDATA encoded field which allows for rich text
	// such as html links: `<a href="http://www.apple.com">Apple</a>`.
	//
	// Use podcast.New(...) to populate this field correctly.
	Description Description `xml:"description"`

	Category       string `xml:"category,omitempty"`
	Cloud          string `xml:"cloud,omitempty"`
	Copyright      string `xml:"copyright,omitempty"`
	Docs           string `xml:"docs,omitempty"`
	Generator      string `xml:"generator,omitempty"`
	Language       string `xml:"language,omitempty"`
	LastBuildDate  string `xml:"lastBuildDate,omitempty"`
	ManagingEditor string `xml:"managingEditor,omitempty"`
	PubDate        string `xml:"pubDate,omitempty"`
	Image          *Image
	Rating         string `xml:"rating,omitempty"`
	SkipHours      string `xml:"skipHours,omitempty"`
	SkipDays       string `xml:"skipDays,omitempty"`
	TTL            int    `xml:"ttl,omitempty"`
	WebMaster      string `xml:"webMaster,omitempty"`
	TextInput      *TextInput
	AtomLink       *AtomLink

	// https://help.apple.com/itc/podcasts_connect/#/itcb54353390
	IAuthor   string `xml:"itunes:author,omitempty"`
	ISubtitle string `xml:"itunes:subtitle,omitempty"`
	ISummary  *ISummary
	IImage    *IImage

	// IExplicit is the parental-advisory flag. Use SetExplicit to populate it.
	IExplicit *explicitFlag `xml:"itunes:explicit,omitempty"`

	// IComplete marks the feed as finished (no future episodes). Use
	// SetComplete to populate it; a false value omits the tag.
	IComplete *yesFlag `xml:"itunes:complete,omitempty"`

	INewFeedURL string `xml:"itunes:new-feed-url,omitempty"`

	// IBlock hides the entire show from the Apple directory. Use SetBlock to
	// populate it; a false value omits the tag.
	IBlock *yesFlag `xml:"itunes:block,omitempty"`

	IDuration   string `xml:"itunes:duration,omitempty"`
	IType       *IType
	IOwner      *Author `xml:"itunes:owner,omitempty"`
	ICategories []*ICategory
	ITitle      string `xml:"itunes:title,omitempty"`

	// Items is a collection of 0..n episodes for this podcast.
	Items  []*Item
	encode func(w io.Writer, o any) error
}

// New instantiates a Podcast with required parameters.
//
// Zero-value time fields default to the current UTC time; non-zero values are
// formatted to the expected proper formats.
func New(title, link, description string,
	pubDate, lastBuildDate time.Time,
) Podcast {
	return Podcast{
		Title:         title,
		Description:   parseDescription(description),
		Link:          link,
		Generator:     fmt.Sprintf("go podcast v%s (github.com/hbmartin/podcast-rss-generator/v2)", pVersion),
		PubDate:       parseDateRFC1123Z(pubDate),
		LastBuildDate: parseDateRFC1123Z(lastBuildDate),
		Language:      "en-us",

		// setup dependency (could inject later)
		encode: encoder,
	}
}

// AddAuthor adds the specified Author to the podcast's ManagingEditor and
// iTunes author tags. When both name and email are supplied, it also sets the
// structured iTunes owner contact.
//
// If email is empty this method is a no-op, since email is the required part
// of an RSS author value.
func (p *Podcast) AddAuthor(name, email string) {
	if len(email) == 0 {
		return
	}
	a := &Author{
		Name:  name,
		Email: email,
	}
	p.ManagingEditor = parseAuthorNameEmail(a)
	p.IAuthor = p.ManagingEditor
	if len(name) != 0 {
		p.IOwner = a
	} else {
		p.IOwner = nil
	}
}

// AddAtomLink adds a FQDN reference to an atom feed.
//
// If href is empty this method is a no-op.
func (p *Podcast) AddAtomLink(href string) {
	if len(href) == 0 {
		return
	}
	p.AtomLink = &AtomLink{
		HREF: href,
		Rel:  "self",
		Type: "application/rss+xml",
	}
}

// AddCategory adds the category to the Podcast.
//
// ICategory can be listed multiple times.
//
// Calling this method multiple times will APPEND the category to the existing
// list, if any, including ICategory.
//
// Note that Apple iTunes has a specific list of categories that only can be
// used and will invalidate the feed if deviated from the list.  The list
// changes occasionally.  Please refer to the following link for the updated
// list:
//
// https://help.apple.com/itc/podcasts_connect/#/itc9267a2f12
//
// If category is empty this method is a no-op. Empty entries in subCategories
// are skipped.
func (p *Podcast) AddCategory(category string, subCategories []string) {
	if len(category) == 0 {
		return
	}

	// RSS 2.0 Category only supports 1-tier
	if len(p.Category) > 0 {
		p.Category = p.Category + "," + category
	} else {
		p.Category = category
	}

	ic := ICategory{
		Text:        category,
		ICategories: make([]*ICategory, 0, len(subCategories)),
	}
	for _, c := range subCategories {
		if len(c) == 0 {
			continue
		}
		ic2 := ICategory{Text: c}
		ic.ICategories = append(ic.ICategories, &ic2)
	}
	p.ICategories = append(p.ICategories, &ic)
}

// AddImage adds the specified Image to the Podcast.
//
// Podcast feeds contain artwork that is a minimum size of
// 1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
// 72 dpi, in JPEG or PNG format with appropriate file
// extensions (.jpg, .png), and in the RGB colorspace. To optimize
// images for mobile devices, Apple recommends compressing your
// image files.
//
// If url is empty this method is a no-op.
func (p *Podcast) AddImage(url string) {
	if len(url) == 0 {
		return
	}
	p.Image = &Image{
		URL:   url,
		Title: p.Title,
		Link:  p.Link,
	}
	p.IImage = &IImage{HREF: url}
}

// AddItem adds the podcast episode.  It returns a count of Items added or any
// errors in validation that may have occurred.
//
// This method takes the "itunes overrides" approach to populating
// itunes tags according to the overrides rules in the specification.
// This not only complies completely with iTunes parsing rules; but, it also
// displays what is possible to be set on an individual episode level – if you
// wish to have more fine grain control over your content.
//
// This method imposes strict validation of the Item being added to confirm
// to Podcast and iTunes specifications.
//
// Article minimal requirements are:
//
//   - Title
//   - Description
//   - Link
//
// Audio, Video and Downloads minimal requirements are:
//
//   - Title
//   - Description
//   - Enclosure (HREF, Type and Length all required)
//
// The following fields are always overwritten (don't set them):
//
//   - GUID
//   - PubDateFormatted
//   - AuthorFormatted
//   - Enclosure.TypeFormatted
//   - Enclosure.LengthFormatted
//
// Recommendations:
//
//   - Just set the minimal fields: the rest get set for you.
//   - Always set an Enclosure.Length, to be nice to your downloaders.
//   - Follow Apple's best practices to enrich your podcasts:
//     https://help.apple.com/itc/podcasts_connect/#/itc2b3780e76
//   - For specifications of itunes tags, see:
//     https://help.apple.com/itc/podcasts_connect/#/itcb54353390
//
// The Item is taken by value and its nested pointer fields (Enclosure, Author,
// IImage, ISummary, IEpisodeType) are deep-cloned before storage. As a result
// the stored episode is a snapshot: mutating the caller's Item, or the structs
// it points at, after AddItem returns does NOT change the feed. Set every
// field before calling AddItem, or re-add the Item to apply later changes.
func (p *Podcast) AddItem(i Item) (int, error) {
	if p == nil {
		return 0, fmt.Errorf("adding item: %w", ErrPodcastRequired)
	}

	i = cloneItem(i)

	if err := validateItem(i); err != nil {
		return len(p.Items), err
	}

	// corrective actions and overrides
	//
	i.PubDateFormatted = parseDateRFC1123Z(i.PubDate)
	i.AuthorFormatted = parseAuthorNameEmail(i.Author)
	if i.Enclosure != nil {
		if len(i.GUID) == 0 {
			i.GUID = i.Enclosure.URL // yep, GUID is the Permlink URL
		}

		if i.Enclosure.Length < 0 {
			i.Enclosure.Length = 0
		}
		i.Enclosure.LengthFormatted = strconv.FormatInt(i.Enclosure.Length, 10)
		i.Enclosure.TypeFormatted = i.Enclosure.Type.String()

		// allow Link to be set for article references to Downloads,
		// otherwise set it to the enclosurer's URL.
		if len(i.Link) == 0 {
			i.Link = i.Enclosure.URL
		}
	} else {
		i.GUID = i.Link // yep, GUID is the Permlink URL
	}

	// iTunes it
	//
	if len(i.IAuthor) == 0 {
		switch {
		case i.Author != nil:
			i.IAuthor = i.Author.Email
		case len(p.IAuthor) != 0:
			i.Author = &Author{Email: p.IAuthor}
			i.IAuthor = p.IAuthor
		case len(p.ManagingEditor) != 0:
			i.Author = &Author{Email: p.ManagingEditor}
			i.IAuthor = p.ManagingEditor
		}
	}
	if i.IImage == nil {
		if p.Image != nil {
			i.IImage = &IImage{HREF: p.Image.URL}
		}
	}

	p.Items = append(p.Items, &i)
	return len(p.Items), nil
}

// AddItems appends multiple episodes in order by calling AddItem for each. It
// stops at the first item that fails validation and returns that error along
// with the number of items successfully stored so far; on success it returns
// the total item count. Each Item follows the same snapshot semantics as
// AddItem.
func (p *Podcast) AddItems(items ...Item) (int, error) {
	if p == nil {
		return 0, fmt.Errorf("adding items: %w", ErrPodcastRequired)
	}
	for _, item := range items {
		if _, err := p.AddItem(item); err != nil {
			return len(p.Items), err
		}
	}
	return len(p.Items), nil
}

// validateItem enforces the required-field rules shared by AddItem and
// Podcast.Validate.
func validateItem(i Item) error {
	if len(i.Title) == 0 || len(i.Description) == 0 {
		return newItemValidationError(i.Title, ErrTitleDescriptionRequired)
	}
	if i.Enclosure != nil {
		if len(i.Enclosure.URL) == 0 {
			return newItemValidationError(i.Title, ErrEnclosureURLRequired)
		}
		if i.Enclosure.Type.String() == enclosureDefault {
			return newItemValidationError(i.Title, ErrEnclosureTypeRequired)
		}
	} else if len(i.Link) == 0 {
		return newItemValidationError(i.Title, ErrLinkRequired)
	}
	return nil
}

// AddPubDate adds the datetime as a parsed PubDate.
//
// UTC time is used by default.
func (p *Podcast) AddPubDate(datetime time.Time) {
	p.PubDate = parseDateRFC1123Z(datetime)
}

// AddLastBuildDate adds the datetime as a parsed PubDate.
//
// UTC time is used by default.
func (p *Podcast) AddLastBuildDate(datetime time.Time) {
	p.LastBuildDate = parseDateRFC1123Z(datetime)
}

// AddSubTitle adds the iTunes subtitle that is displayed with the title
// in iTunes.
//
// Note that this field should be just a few words long according to Apple.
// This method will truncate the string to 64 chars if too long with "...".
//
// If subTitle is empty this method is a no-op.
func (p *Podcast) AddSubTitle(subTitle string) {
	if len(subTitle) == 0 {
		return
	}
	p.ISubtitle = truncateRunesWithSuffix(subTitle, iTunesSubtitleRuneLimit, "...")
}

// AddSummary adds the iTunes summary.
//
// Limit: 4000 characters
//
// Note that this field is a CDATA encoded field which allows for rich text
// such as html links: `<a href="http://www.apple.com">Apple</a>`.
//
// If summary is empty this method is a no-op.
func (p *Podcast) AddSummary(summary string) {
	if len(summary) == 0 {
		return
	}
	p.ISummary = &ISummary{
		Text: truncateRunes(summary, iTunesSummaryRuneLimit),
	}
}

// Bytes returns an encoded []byte slice.
func (p *Podcast) Bytes() []byte {
	return []byte(p.String())
}

// Encode writes the bytes to the io.Writer stream in RSS 2.0 specification.
func (p *Podcast) Encode(w io.Writer) error {
	if p == nil {
		return fmt.Errorf("encoding podcast: %w", ErrPodcastRequired)
	}
	if isNilWriter(w) {
		return fmt.Errorf("encoding podcast: %w", ErrWriterRequired)
	}
	if _, err := w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")); err != nil {
		return fmt.Errorf("writing xml header: %w", err)
	}

	atomLink := ""
	if p.AtomLink != nil {
		atomLink = "http://www.w3.org/2005/Atom"
	}
	contentNS := ""
	for _, it := range p.Items {
		if it != nil && it.ContentEncoded != nil {
			contentNS = contentNamespace
			break
		}
	}
	wrapped := podcastWrapper{
		ITUNESNS:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
		ATOMNS:    atomLink,
		CONTENTNS: contentNS,
		Version:   "2.0",
		Channel:   p,
	}
	encode := p.encode
	if encode == nil {
		encode = encoder
	}
	return encode(w, wrapped)
}

// String encodes the Podcast state to a string.
func (p *Podcast) String() string {
	b := new(bytes.Buffer)
	if err := p.Encode(b); err != nil {
		return "String: podcast.write returned the error: " + err.Error()
	}
	return b.String()
}

// // Write implements the io.Writer interface to write an RSS 2.0 stream
// // that is compliant to the RSS 2.0 specification.
// func (p *Podcast) Write(b []byte) (n int, err error) {
// 	buf := bytes.NewBuffer(b)
// 	if err := p.Encode(buf); err != nil {
// 		return 0, fmt.Errorf("Write: podcast.encode returned error: %w", err)
// 	}
// 	return buf.Len(), nil
// }

type podcastWrapper struct {
	XMLName   xml.Name `xml:"rss"`
	Version   string   `xml:"version,attr"`
	ATOMNS    string   `xml:"xmlns:atom,attr,omitempty"`
	ITUNESNS  string   `xml:"xmlns:itunes,attr"`
	CONTENTNS string   `xml:"xmlns:content,attr,omitempty"`
	Channel   *Podcast
}

func encoder(w io.Writer, o any) error {
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	if err := e.Encode(o); err != nil {
		return fmt.Errorf("encoding xml: %w", err)
	}
	return nil
}

func parseAuthorNameEmail(a *Author) string {
	var author string
	if a != nil {
		author = a.Email
		if len(a.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", a.Email, a.Name)
		}
	}
	return author
}

// AddType adds the Apple Podcasts show type.
func (p *Podcast) AddType(podcastType PodcastType) {
	p.IType = &IType{Text: podcastType.String()}
}

// SetExplicit sets the itunes:explicit parental-advisory flag for the show. A
// true value renders "true" (the Explicit badge); a false value renders
// "false" (the Clean badge). Both are serialized.
func (p *Podcast) SetExplicit(explicit bool) {
	flag := explicitFlag(explicit)
	p.IExplicit = &flag
}

// SetComplete marks whether the show is complete, meaning no further episodes
// will be published. A true value renders "Yes"; a false value omits the tag,
// which Apple treats as "not complete".
func (p *Podcast) SetComplete(complete bool) {
	if !complete {
		p.IComplete = nil
		return
	}
	flag := yesFlag(true)
	p.IComplete = &flag
}

// SetBlock sets whether the entire show is hidden from the Apple directory. A
// true value renders "Yes"; a false value omits the tag, which Apple treats as
// "not blocked".
func (p *Podcast) SetBlock(block bool) {
	if !block {
		p.IBlock = nil
		return
	}
	flag := yesFlag(true)
	p.IBlock = &flag
}

// Validate reports whether the podcast satisfies the channel- and item-level
// requirements enforced by this package. It returns nil when the feed is
// valid, or a joined error aggregating every problem found so callers can
// surface all of them at once.
//
// Channel requirements are Title, Description and Link. Each item is validated
// with the same rules AddItem applies, so a Podcast assembled entirely through
// AddItem only needs its channel fields checked here.
func (p *Podcast) Validate() error {
	if p == nil {
		return fmt.Errorf("validating podcast: %w", ErrPodcastRequired)
	}

	var errs []error
	if len(p.Title) == 0 {
		errs = append(errs, ErrChannelTitleRequired)
	}
	if len(p.Description) == 0 {
		errs = append(errs, ErrChannelDescriptionRequired)
	}
	if len(p.Link) == 0 {
		errs = append(errs, ErrChannelLinkRequired)
	}
	for _, it := range p.Items {
		if it == nil {
			continue
		}
		if err := validateItem(*it); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// AddChannelType adds the Apple Podcasts show type from a string.
//
// Deprecated: use AddType with podcast.Episodic or podcast.Serial.
func (p *Podcast) AddChannelType(channelType string) {
	if parsed, ok := parseType(channelType); ok {
		p.IType = &IType{Text: parsed.String()}
	}
}

func parseDescription(d string) Description {
	if len(d) <= podcastDescriptionByteLimit {
		return Description(d)
	}

	limit := podcastDescriptionByteLimit
	for limit > 0 && !utf8.RuneStart(d[limit]) {
		limit--
	}
	return Description(d[:limit])
}

func parseType(channelType string) (PodcastType, bool) {
	switch strings.ToLower(strings.TrimSpace(channelType)) {
	case Episodic.String():
		return Episodic, true
	case Serial.String():
		return Serial, true
	}
	return Episodic, false
}

func newItemValidationError(title string, err error) error {
	return &ItemValidationError{
		Title: title,
		Err:   err,
	}
}

func cloneItem(i Item) Item {
	if i.Enclosure != nil {
		enclosure := *i.Enclosure
		i.Enclosure = &enclosure
	}
	if i.Author != nil {
		author := *i.Author
		i.Author = &author
	}
	if i.IImage != nil {
		image := *i.IImage
		i.IImage = &image
	}
	if i.ContentEncoded != nil {
		content := *i.ContentEncoded
		i.ContentEncoded = &content
	}
	if i.ISummary != nil {
		summary := *i.ISummary
		i.ISummary = &summary
	}
	if i.IEpisodeType != nil {
		episodeType := *i.IEpisodeType
		i.IEpisodeType = &episodeType
	}
	if i.IExplicit != nil {
		explicit := *i.IExplicit
		i.IExplicit = &explicit
	}
	if i.IBlock != nil {
		block := *i.IBlock
		i.IBlock = &block
	}
	return i
}

func isNilWriter(w io.Writer) bool {
	if w == nil {
		return true
	}
	v := reflect.ValueOf(w)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func truncateRunes(s string, limit int) string {
	if limit <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= limit {
		return s
	}
	return string([]rune(s)[:limit])
}

func truncateRunesWithSuffix(s string, limit int, suffix string) string {
	if utf8.RuneCountInString(s) <= limit {
		return s
	}
	suffixLen := utf8.RuneCountInString(suffix)
	if suffixLen >= limit {
		return truncateRunes(suffix, limit)
	}
	return truncateRunes(s, limit-suffixLen) + suffix
}
