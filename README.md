[![Go Reference](https://pkg.go.dev/badge/github.com/hbmartin/podcast-rss-generator/v2.svg)](https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2)
[![CI](https://github.com/hbmartin/podcast-rss-generator/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/hbmartin/podcast-rss-generator/actions/workflows/ci.yml)
[![Coverage Status](https://codecov.io/gh/hbmartin/podcast-rss-generator/branch/master/graph/badge.svg)](https://codecov.io/gh/hbmartin/podcast-rss-generator)
[![Go Report Card](https://goreportcard.com/badge/github.com/hbmartin/podcast-rss-generator/v2)](https://goreportcard.com/report/github.com/hbmartin/podcast-rss-generator/v2)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)



# podcast
`import "github.com/hbmartin/podcast-rss-generator/v2"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package podcast generates RSS 2.0 podcast feeds with common Apple Podcasts
tags using a small Go API.

Full documentation with detailed examples is located at <a href="https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2">https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2</a>

# Usage

To use, `go get` and `import` the package like your typical Go library.


	$ go get github.com/hbmartin/podcast-rss-generator/v2@latest
	
	import "github.com/hbmartin/podcast-rss-generator/v2"

The API exposes method receivers on feed structs that fill derived RSS and
Apple Podcasts fields, validate required item fields, and keep generated XML
consistent.

Notably, the `Podcast.AddItem` function performs most
of the heavy lifting by taking the `Item` input and performing
validation, overrides and duplicate setters through the feed.

Detailed examples of the API are at <a href="https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2">https://pkg.go.dev/github.com/hbmartin/podcast-rss-generator/v2</a>.

# Contributing

See the CONTRIBUTING.md for all the details.

# Go Modules

This library is supported on Go 1.24.0 and higher.

The repository uses Go modules and a mise-managed toolchain for local and CI
checks.

# Extensibility

The exported structs remain available for callers that need direct control
over RSS and Apple Podcasts fields. Prefer the Add methods for fields with
formatting, validation, or derived values.

# Fuzzing

Native Go fuzzing covers feed encoding with XML-sensitive text, invalid UTF-8,
long descriptions, enclosure metadata, and date bytes.


	mise run fuzz-smoke

To run longer fuzzing locally:


	go test -run '^$' -fuzz=FuzzPodcastEncode -fuzztime=1m ./...

# Roadmap

Current v2 work focuses on preserving the public API, improving validation
and safety, and keeping generated documentation and examples accurate.

# Versioning

Releases follow semantic versioning. Public API removals or incompatible
behavior changes require a new major version.

# Release Notes

See CHANGELOG.md in the repository for release notes. The changelog is the
source of truth for published version history.

# References

RSS 2.0: <a href="https://cyber.harvard.edu/rss/rss.html">https://cyber.harvard.edu/rss/rss.html</a>

Podcasts: <a href="https://help.apple.com/itc/podcasts_connect/#/itca5b22233">https://help.apple.com/itc/podcasts_connect/#/itca5b22233</a>




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type AtomLink](#AtomLink)
* [type Author](#Author)
* [type Description](#Description)
  * [func (d Description) MarshalXML(e *xml.Encoder, start xml.StartElement) error](#Description.MarshalXML)
* [type Enclosure](#Enclosure)
* [type EnclosureType](#EnclosureType)
  * [func (et EnclosureType) String() string](#EnclosureType.String)
* [type EpisodeType](#EpisodeType)
  * [func (et EpisodeType) String() string](#EpisodeType.String)
* [type ICategory](#ICategory)
* [type IEpisodeType](#IEpisodeType)
* [type IImage](#IImage)
* [type ISummary](#ISummary)
* [type IType](#IType)
* [type Image](#Image)
* [type Item](#Item)
  * [func (i *Item) AddDescription(d string)](#Item.AddDescription)
  * [func (i *Item) AddDuration(durationInSeconds int64)](#Item.AddDuration)
  * [func (i *Item) AddEnclosure(url string, enclosureType EnclosureType, lengthInBytes int64)](#Item.AddEnclosure)
  * [func (i *Item) AddEpisode()](#Item.AddEpisode)
  * [func (i *Item) AddEpisodeType(episodeType EpisodeType)](#Item.AddEpisodeType)
  * [func (i *Item) AddImage(url string)](#Item.AddImage)
  * [func (i *Item) AddPubDate(datetime time.Time)](#Item.AddPubDate)
  * [func (i *Item) AddSummary(summary string)](#Item.AddSummary)
* [type ItemValidationError](#ItemValidationError)
  * [func (e *ItemValidationError) Error() string](#ItemValidationError.Error)
  * [func (e *ItemValidationError) Unwrap() error](#ItemValidationError.Unwrap)
* [type Podcast](#Podcast)
  * [func New(title, link, description string, pubDate, lastBuildDate time.Time) Podcast](#New)
  * [func (p *Podcast) AddAtomLink(href string)](#Podcast.AddAtomLink)
  * [func (p *Podcast) AddAuthor(name, email string)](#Podcast.AddAuthor)
  * [func (p *Podcast) AddCategory(category string, subCategories []string)](#Podcast.AddCategory)
  * [func (p *Podcast) AddChannelType(channelType string)](#Podcast.AddChannelType)
  * [func (p *Podcast) AddImage(url string)](#Podcast.AddImage)
  * [func (p *Podcast) AddItem(i Item) (int, error)](#Podcast.AddItem)
  * [func (p *Podcast) AddLastBuildDate(datetime time.Time)](#Podcast.AddLastBuildDate)
  * [func (p *Podcast) AddPubDate(datetime time.Time)](#Podcast.AddPubDate)
  * [func (p *Podcast) AddSubTitle(subTitle string)](#Podcast.AddSubTitle)
  * [func (p *Podcast) AddSummary(summary string)](#Podcast.AddSummary)
  * [func (p *Podcast) AddType(podcastType PodcastType)](#Podcast.AddType)
  * [func (p *Podcast) Bytes() []byte](#Podcast.Bytes)
  * [func (p *Podcast) Encode(w io.Writer) error](#Podcast.Encode)
  * [func (p *Podcast) String() string](#Podcast.String)
* [type PodcastType](#PodcastType)
  * [func (pt PodcastType) String() string](#PodcastType.String)
* [type TextInput](#TextInput)

#### <a name="pkg-examples">Examples</a>
* [Item.AddDuration](#example_Item_AddDuration)
* [Item.AddPubDate](#example_Item_AddPubDate)
* [New](#example_New)
* [Podcast.AddAuthor](#example_Podcast_AddAuthor)
* [Podcast.AddCategory](#example_Podcast_AddCategory)
* [Podcast.AddImage](#example_Podcast_AddImage)
* [Podcast.AddItem](#example_Podcast_AddItem)
* [Podcast.AddLastBuildDate](#example_Podcast_AddLastBuildDate)
* [Podcast.AddPubDate](#example_Podcast_AddPubDate)
* [Podcast.AddSummary](#example_Podcast_AddSummary)
* [Podcast.Bytes](#example_Podcast_Bytes)
* [Package (HttpHandlers)](#example__httpHandlers)
* [Package (IoWriter)](#example__ioWriter)


## <a name="pkg-variables">Variables</a>
``` go
var (
    ErrPodcastRequired          = errors.New("podcast is required")
    ErrWriterRequired           = errors.New("writer is required")
    ErrTitleDescriptionRequired = errors.New("title and description are required")
    ErrEnclosureURLRequired     = errors.New("enclosure url is required")
    ErrEnclosureTypeRequired    = errors.New("enclosure type is required")
    ErrLinkRequired             = errors.New("link is required when not using enclosure")
)
```
Sentinel errors returned by feed and item validation.




## <a name="AtomLink">type</a> [AtomLink](./atomlink.go#L6-L11)
``` go
type AtomLink struct {
    XMLName xml.Name `xml:"atom:link"`
    HREF    string   `xml:"href,attr"`
    Rel     string   `xml:"rel,attr"`
    Type    string   `xml:"type,attr"`
}

```
AtomLink represents the Atom reference link.










## <a name="Author">type</a> [Author](./author.go#L8-L12)
``` go
type Author struct {
    XMLName xml.Name
    Name    string `xml:"itunes:name"`
    Email   string `xml:"itunes:email"`
}

```
Author represents a named author and email.

For iTunes compliance, both Name and Email are required.










## <a name="Description">type</a> [Description](./description.go#L11)
``` go
type Description string
```
Description is a rich-text field for channel and episode description tags.

This is rendered as CDATA which allows for HTML tags such as `<a href="">`.










### <a name="Description.MarshalXML">func</a> (Description) [MarshalXML](./description.go#L15)
``` go
func (d Description) MarshalXML(e *xml.Encoder, start xml.StartElement) error
```
MarshalXML renders the description text as CDATA so rich text such as
`<a href="">` stays valid XML.




## <a name="Enclosure">type</a> [Enclosure](./enclosure.go#L47-L66)
``` go
type Enclosure struct {
    XMLName xml.Name `xml:"enclosure"`

    // URL is the downloadable url for the content. (Required)
    URL string `xml:"url,attr"`

    // Length is the size in Bytes of the download. (Required)
    Length int64 `xml:"-"`
    // LengthFormatted is the size in Bytes of the download. (Required)
    //
    // This field gets overwritten with the API when setting Length.
    LengthFormatted string `xml:"length,attr"`

    // Type is MIME type encoding of the download. (Required)
    Type EnclosureType `xml:"-"`
    // TypeFormatted is MIME type encoding of the download. (Required)
    //
    // This field gets overwritten with the API when setting Type.
    TypeFormatted string `xml:"type,attr"`
}

```
Enclosure represents a download enclosure.










## <a name="EnclosureType">type</a> [EnclosureType](./enclosure.go#L22)
``` go
type EnclosureType int
```
EnclosureType specifies the type of the enclosure.


``` go
const (
    EnclosureUnknown EnclosureType = iota
    M4A
    M4V
    MP4
    MP3
    MOV
    PDF
    EPUB
)
```
EnclosureType specifies the type of the enclosure.










### <a name="EnclosureType.String">func</a> (EnclosureType) [String](./enclosure.go#L25)
``` go
func (et EnclosureType) String() string
```
String returns the MIME type encoding of the specified EnclosureType.




## <a name="EpisodeType">type</a> [EpisodeType](./type.go#L62)
``` go
type EpisodeType int
```
EpisodeType specifies whether an episode is full, trailer, or bonus content.


``` go
const (
    Full EpisodeType = iota
    Trailer
    Bonus
)
```
EpisodeType specifies the type of an episode.










### <a name="EpisodeType.String">func</a> (EpisodeType) [String](./type.go#L65)
``` go
func (et EpisodeType) String() string
```
String returns the Apple Podcasts encoding of the specified EpisodeType.




## <a name="ICategory">type</a> [ICategory](./itunes.go#L9-L13)
``` go
type ICategory struct {
    XMLName     xml.Name `xml:"itunes:category"`
    Text        string   `xml:"text,attr"`
    ICategories []*ICategory
}

```
ICategory is a 2-tier classification system for iTunes.










## <a name="IEpisodeType">type</a> [IEpisodeType](./itunes.go#L43-L46)
``` go
type IEpisodeType struct {
    XMLName xml.Name `xml:"itunes:episodeType"`
    Text    string   `xml:",chardata"`
}

```
IEpisodeType renders podcast.EpisodeType.










## <a name="IImage">type</a> [IImage](./itunes.go#L23-L26)
``` go
type IImage struct {
    XMLName xml.Name `xml:"itunes:image"`
    HREF    string   `xml:"href,attr"`
}

```
IImage represents an iTunes image.

Podcast feeds contain artwork that is a minimum size of
1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
72 dpi, in JPEG or PNG format with appropriate file
extensions (.jpg, .png), and in the RGB colorspace. To optimize
images for mobile devices, Apple recommends compressing your
image files.










## <a name="ISummary">type</a> [ISummary](./itunes.go#L31-L34)
``` go
type ISummary struct {
    XMLName xml.Name `xml:"itunes:summary"`
    Text    string   `xml:",cdata"`
}

```
ISummary is a 4000 character rich-text field for the itunes:summary tag.

This is rendered as CDATA which allows for HTML tags such as `<a href="">`.










## <a name="IType">type</a> [IType](./itunes.go#L37-L40)
``` go
type IType struct {
    XMLName xml.Name `xml:"itunes:type"`
    Text    string   `xml:",chardata"`
}

```
IType renders podcast.PodcastType.










## <a name="Image">type</a> [Image](./image.go#L13-L21)
``` go
type Image struct {
    XMLName     xml.Name `xml:"image"`
    URL         string   `xml:"url"`
    Title       string   `xml:"title"`
    Link        string   `xml:"link"`
    Description string   `xml:"description,omitempty"`
    Width       int      `xml:"width,omitempty"`
    Height      int      `xml:"height,omitempty"`
}

```
Image represents an image.

Podcast feeds contain artwork that is a minimum size of
1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
72 dpi, in JPEG or PNG format with appropriate file
extensions (.jpg, .png), and in the RGB colorspace. To optimize
images for mobile devices, Apple recommends compressing your
image files.










## <a name="Item">type</a> [Item](./item.go#L26-L170)
``` go
type Item struct {
    XMLName xml.Name `xml:"item"`

    // GUID is the episode’s globally unique identifier (GUID). It is
    // recommended to set this tag on each item.
    //
    // It is very important that each episode have a unique GUID and
    // that it never changes, even if an episode’s metadata, like
    // title or enclosure URL, do change.
    GUID string `xml:"guid"`

    // Title is an episode title. It is a required field per iTunes
    // definitions.
    Title string `xml:"title"`

    // Link is an episode link URL.
    //
    // Do not use HTMl tags.  Only raw URLs such as https:// are allowed.
    Link string `xml:"link"`

    // Description is text containing one or more sentences describing
    // your episode to potential listeners.
    //
    // Use item.AddDescription(...) to populate this field correctly.
    Description Description `xml:"description"`

    // PubDate is the date and time when an episode was released. It is
    // recommended to set this tag on each item.
    //
    // Use item.AddPubDate(...) to populate this field correctly.
    PubDate time.Time `xml:"-"`

    // PubDateFormatted is deprecated.  Do not populate nor read this
    // string as it will be removed in a future release.
    PubDateFormatted string `xml:"pubDate,omitempty"`

    // Enclosure is of type podcast.Enclosure. It is a required field per
    // iTunes definitions.
    //
    // Use item.AddEnclosure(...) to populate this field correctly.
    Enclosure *Enclosure

    // IDuration is the duration of an episode.
    //
    // Use item.AddDuration(...) to populate this field correctly.
    IDuration string `xml:"itunes:duration,omitempty"`

    // IExplicit defines the episode parental advisory information.
    //
    // Where the explicit value can be one of the following:
    //
    // "true" : If you specify true, indicating the presence of explicit
    // content, Apple Podcasts displays an Explicit parental advisory
    // graphic for your episode.
    // Episodes containing explicit material aren’t available in some Apple
    // Podcasts territories.
    //
    // "false" : If you specify false, indicating that the episode does not
    // contain explicit language or adult content, Apple Podcasts displays
    // a Clean parental advisory graphic for your episode.
    IExplicit string `xml:"itunes:explicit,omitempty"`

    // ITitle is a Situational episode title specific for Apple Podcasts.
    //
    // This tag is a string containing a clear concise name of your
    // episode on Apple Podcasts.
    //
    // Don’t specify the episode number or season number in this tag. Instead,
    // specify those details in the appropriate tags IEpisode and ISeason>.
    //
    // Also, don’t repeat the title of your show within your episode title.
    //
    // Separating episode and season number from the title makes it possible
    // for Apple to easily index and order content from all shows.
    ITitle string `xml:"itunes:title,omitempty"`

    // IEpisode is a Situational tag for the episode number.
    //
    // If all your episodes have numbers and you would like them to be ordered
    // based on them use this tag for each one.
    //
    // Episode numbers are optional for type episodic shows, but are
    // mandatory for serial shows.
    //
    // Where episode is a non-zero integer (1, 2, 3, etc.) representing your
    // episode number.
    IEpisode string `xml:"itunes:episode,omitempty"`

    // ISeason is a Situational tag for the episode season number.
    //
    // If an episode is within a season use this tag.
    //
    // Where season is a non-zero integer (1, 2, 3, etc.) representing your
    // season number.
    //
    // To allow the season feature for shows containing a single season, if
    // only one season exists in the RSS feed, Apple Podcasts doesn’t display
    // a season number. When you add a second season to the RSS feed, Apple
    // Podcasts displays the season numbers.
    ISeason string `xml:"itunes:season,omitempty"`

    // IEpisodeType is a Situational tag for the episode type.
    //
    // If an episode is a trailer or bonus content, use this tag.
    //
    // Use AddEpisodeType(...) to populate this field correctly.
    IEpisodeType *IEpisodeType

    // IBlock is a Situational tag to show or hide the status of the episode.
    //
    // If you want an episode removed from the Apple directory, use this tag.
    //
    // Specifying the tag with a "Yes" value prevents that episode from
    // appearing in Apple Podcasts.
    //
    // For example, you might want to block a specific episode if you know
    // that its content would otherwise cause the entire podcast to be
    // removed from Apple Podcasts.
    //
    // Specifying any value other than Yes has no effect.
    IBlock string `xml:"itunes:block,omitempty"`

    Author             *Author `xml:"-"`
    AuthorFormatted    string  `xml:"author,omitempty"`
    Category           string  `xml:"category,omitempty"`
    Comments           string  `xml:"comments,omitempty"`
    IAuthor            string  `xml:"itunes:author,omitempty"`
    IIsClosedCaptioned string  `xml:"itunes:isClosedCaptioned,omitempty"`
    IOrder             string  `xml:"itunes:order,omitempty"`
    ISubtitle          string  `xml:"itunes:subtitle,omitempty"`
    ISummary           *ISummary
    IImage             *IImage
    Source             string `xml:"source,omitempty"`
}

```
Item represents a single entry in a podcast.

Article minimal requirements are:
- Title
- Description
- Link

Audio minimal requirements are:
- Title
- Description
- Enclosure (HREF, Type and Length all required)

Recommendations:
- Setting the minimal fields sets most of other fields, including iTunes.
- Use the Published time.Time setting instead of PubDate.
- Always set an Enclosure.Length, to be nice to your downloaders.
- Use Enclosure.Type instead of setting TypeFormatted for valid extensions.










### <a name="Item.AddDescription">func</a> (\*Item) [AddDescription](./item.go#L266)
``` go
func (i *Item) AddDescription(d string)
```
AddDescription adds a rich-text description tag.

Limit: 10000 characters

Note that this field is a CDATA encoded field which allows for rich text
such as html links: `<a href="<a href="http://www.apple.com">http://www.apple.com</a>">Apple</a>`.




### <a name="Item.AddDuration">func</a> (\*Item) [AddDuration](./item.go#L271)
``` go
func (i *Item) AddDuration(durationInSeconds int64)
```
AddDuration adds the duration to the iTunes duration field.




### <a name="Item.AddEnclosure">func</a> (\*Item) [AddEnclosure](./item.go#L173-L174)
``` go
func (i *Item) AddEnclosure(
    url string, enclosureType EnclosureType, lengthInBytes int64)
```
AddEnclosure adds the downloadable asset to the podcast Item.




### <a name="Item.AddEpisode">func</a> (\*Item) [AddEpisode](./item.go#L185)
``` go
func (i *Item) AddEpisode()
```
AddEpisode adds the situational tags with rules to iTunes' episodes.
Using this function will ensure a properly formatted episode has been
added to the feed in compliance to iTunes' requirements.




### <a name="Item.AddEpisodeType">func</a> (\*Item) [AddEpisodeType](./item.go#L190)
``` go
func (i *Item) AddEpisodeType(episodeType EpisodeType)
```
AddEpisodeType adds the Apple Podcasts episode type.




### <a name="Item.AddImage">func</a> (\*Item) [AddImage](./item.go#L234)
``` go
func (i *Item) AddImage(url string)
```
AddImage adds the image as an iTunes-only IImage.  RSS 2.0 does not have
the specification of Images at the Item level.

Podcast feeds contain artwork that is a minimum size of
1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
72 dpi, in JPEG or PNG format with appropriate file
extensions (.jpg, .png), and in the RGB colorspace. To optimize
images for mobile devices, Apple recommends compressing your
image files.




### <a name="Item.AddPubDate">func</a> (\*Item) [AddPubDate](./item.go#L243)
``` go
func (i *Item) AddPubDate(datetime time.Time)
```
AddPubDate adds the datetime as a parsed PubDate.

UTC time is used by default.




### <a name="Item.AddSummary">func</a> (\*Item) [AddSummary](./item.go#L254)
``` go
func (i *Item) AddSummary(summary string)
```
AddSummary adds the iTunes summary.

Limit: 4000 characters

Note that this field is a CDATA encoded field which allows for rich text
such as html links: `<a href="<a href="http://www.apple.com">http://www.apple.com</a>">Apple</a>`.




## <a name="ItemValidationError">type</a> [ItemValidationError](./podcast.go#L31-L34)
``` go
type ItemValidationError struct {
    Title string
    Err   error
}

```
ItemValidationError describes why an item could not be added to a podcast.










### <a name="ItemValidationError.Error">func</a> (\*ItemValidationError) [Error](./podcast.go#L37)
``` go
func (e *ItemValidationError) Error() string
```
Error returns a human-readable item validation message.




### <a name="ItemValidationError.Unwrap">func</a> (\*ItemValidationError) [Unwrap](./podcast.go#L51)
``` go
func (e *ItemValidationError) Unwrap() error
```
Unwrap returns the underlying validation error.




## <a name="Podcast">type</a> [Podcast](./podcast.go#L59-L131)
``` go
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
    IAuthor     string `xml:"itunes:author,omitempty"`
    ISubtitle   string `xml:"itunes:subtitle,omitempty"`
    ISummary    *ISummary
    IImage      *IImage
    IExplicit   string `xml:"itunes:explicit,omitempty"`
    IComplete   string `xml:"itunes:complete,omitempty"`
    INewFeedURL string `xml:"itunes:new-feed-url,omitempty"`
    IBlock      string `xml:"itunes:block,omitempty"`
    IDuration   string `xml:"itunes:duration,omitempty"`
    IType       *IType
    IOwner      *Author `xml:"itunes:owner,omitempty"`
    ICategories []*ICategory
    ITitle      string `xml:"itunes:title,omitempty"`

    // Items is a collection of 0..n episodes for this podcast.
    Items []*Item
    // contains filtered or unexported fields
}

```
Podcast represents a podcast.







### <a name="New">func</a> [New](./podcast.go#L137-L138)
``` go
func New(title, link, description string,
    pubDate, lastBuildDate time.Time) Podcast
```
New instantiates a Podcast with required parameters.

Zero-value time fields default to the current UTC time; non-zero values are
formatted to the expected proper formats.





### <a name="Podcast.AddAtomLink">func</a> (\*Podcast) [AddAtomLink](./podcast.go#L174)
``` go
func (p *Podcast) AddAtomLink(href string)
```
AddAtomLink adds a FQDN reference to an atom feed.




### <a name="Podcast.AddAuthor">func</a> (\*Podcast) [AddAuthor](./podcast.go#L156)
``` go
func (p *Podcast) AddAuthor(name, email string)
```
AddAuthor adds the specified Author to the podcast's ManagingEditor and
iTunes author tags. When both name and email are supplied, it also sets the
structured iTunes owner contact.




### <a name="Podcast.AddCategory">func</a> (\*Podcast) [AddCategory](./podcast.go#L198)
``` go
func (p *Podcast) AddCategory(category string, subCategories []string)
```
AddCategory adds the category to the Podcast.

ICategory can be listed multiple times.

Calling this method multiple times will APPEND the category to the existing
list, if any, including ICategory.

Note that Apple iTunes has a specific list of categories that only can be
used and will invalidate the feed if deviated from the list.  The list
changes occasionally.  Please refer to the following link for the updated
list:

<a href="https://help.apple.com/itc/podcasts_connect/#/itc9267a2f12">https://help.apple.com/itc/podcasts_connect/#/itc9267a2f12</a>




### <a name="Podcast.AddChannelType">func</a> (\*Podcast) [AddChannelType](./podcast.go#L484)
``` go
func (p *Podcast) AddChannelType(channelType string)
```
AddChannelType adds the Apple Podcasts show type from a string.

Deprecated: use AddType with podcast.Episodic or podcast.Serial.




### <a name="Podcast.AddImage">func</a> (\*Podcast) [AddImage](./podcast.go#L232)
``` go
func (p *Podcast) AddImage(url string)
```
AddImage adds the specified Image to the Podcast.

Podcast feeds contain artwork that is a minimum size of
1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
72 dpi, in JPEG or PNG format with appropriate file
extensions (.jpg, .png), and in the RGB colorspace. To optimize
images for mobile devices, Apple recommends compressing your
image files.




### <a name="Podcast.AddItem">func</a> (\*Podcast) [AddItem](./podcast.go#L284)
``` go
func (p *Podcast) AddItem(i Item) (int, error)
```
AddItem adds the podcast episode.  It returns a count of Items added or any
errors in validation that may have occurred.

This method takes the "itunes overrides" approach to populating
itunes tags according to the overrides rules in the specification.
This not only complies completely with iTunes parsing rules; but, it also
displays what is possible to be set on an individual episode level – if you
wish to have more fine grain control over your content.

This method imposes strict validation of the Item being added to confirm
to Podcast and iTunes specifications.

Article minimal requirements are:


	- Title
	- Description
	- Link

Audio, Video and Downloads minimal requirements are:


	- Title
	- Description
	- Enclosure (HREF, Type and Length all required)

The following fields are always overwritten (don't set them):


	- GUID
	- PubDateFormatted
	- AuthorFormatted
	- Enclosure.TypeFormatted
	- Enclosure.LengthFormatted

Recommendations:


	- Just set the minimal fields: the rest get set for you.
	- Always set an Enclosure.Length, to be nice to your downloaders.
	- Follow Apple's best practices to enrich your podcasts:
	  <a href="https://help.apple.com/itc/podcasts_connect/#/itc2b3780e76">https://help.apple.com/itc/podcasts_connect/#/itc2b3780e76</a>
	- For specifications of itunes tags, see:
	  <a href="https://help.apple.com/itc/podcasts_connect/#/itcb54353390">https://help.apple.com/itc/podcasts_connect/#/itcb54353390</a>




### <a name="Podcast.AddLastBuildDate">func</a> (\*Podcast) [AddLastBuildDate](./podcast.go#L364)
``` go
func (p *Podcast) AddLastBuildDate(datetime time.Time)
```
AddLastBuildDate adds the datetime as a parsed PubDate.

UTC time is used by default.




### <a name="Podcast.AddPubDate">func</a> (\*Podcast) [AddPubDate](./podcast.go#L357)
``` go
func (p *Podcast) AddPubDate(datetime time.Time)
```
AddPubDate adds the datetime as a parsed PubDate.

UTC time is used by default.




### <a name="Podcast.AddSubTitle">func</a> (\*Podcast) [AddSubTitle](./podcast.go#L373)
``` go
func (p *Podcast) AddSubTitle(subTitle string)
```
AddSubTitle adds the iTunes subtitle that is displayed with the title
in iTunes.

Note that this field should be just a few words long according to Apple.
This method will truncate the string to 64 chars if too long with "...".




### <a name="Podcast.AddSummary">func</a> (\*Podcast) [AddSummary](./podcast.go#L386)
``` go
func (p *Podcast) AddSummary(summary string)
```
AddSummary adds the iTunes summary.

Limit: 4000 characters

Note that this field is a CDATA encoded field which allows for rich text
such as html links: `<a href="<a href="http://www.apple.com">http://www.apple.com</a>">Apple</a>`.




### <a name="Podcast.AddType">func</a> (\*Podcast) [AddType](./podcast.go#L477)
``` go
func (p *Podcast) AddType(podcastType PodcastType)
```
AddType adds the Apple Podcasts show type.




### <a name="Podcast.Bytes">func</a> (\*Podcast) [Bytes](./podcast.go#L396)
``` go
func (p *Podcast) Bytes() []byte
```
Bytes returns an encoded []byte slice.




### <a name="Podcast.Encode">func</a> (\*Podcast) [Encode](./podcast.go#L401)
``` go
func (p *Podcast) Encode(w io.Writer) error
```
Encode writes the bytes to the io.Writer stream in RSS 2.0 specification.




### <a name="Podcast.String">func</a> (\*Podcast) [String](./podcast.go#L430)
``` go
func (p *Podcast) String() string
```
String encodes the Podcast state to a string.




## <a name="PodcastType">type</a> [PodcastType](./type.go#L36)
``` go
type PodcastType int
```
PodcastType specifies the type of the podcast.

Its values can be one of the following:

Episodic (default). Specify episodic when episodes are intended to be
consumed without any specific order. Apple Podcasts will present newest
episodes first and display the publish date (required) of each episode.
If organized into seasons, the newest season will be presented first
- otherwise, episodes will be grouped by year published, newest first.

For new subscribers, Apple Podcasts adds the newest, most recent episode
in their Library.

Serial. Specify serial when episodes are intended to be consumed in
sequential order. Apple Podcasts will present the oldest episodes
first and display the episode numbers (required) of each episode. If
organized into seasons, the newest season will be presented first and
<itunes:episode> numbers must be given for each episode.

For new subscribers, Apple Podcasts adds the first episode to their
Library, or the entire current season if using seasons.


``` go
const (
    Episodic PodcastType = iota
    Serial
)
```
Episodic and Serial are the supported PodcastType values.










### <a name="PodcastType.String">func</a> (PodcastType) [String](./type.go#L39)
``` go
func (pt PodcastType) String() string
```
String returns the Apple Podcasts encoding of the specified PodcastType.




## <a name="TextInput">type</a> [TextInput](./textinput.go#L6-L12)
``` go
type TextInput struct {
    XMLName     xml.Name `xml:"textInput"`
    Title       string   `xml:"title"`
    Description string   `xml:"description"`
    Name        string   `xml:"name"`
    Link        string   `xml:"link"`
}

```
TextInput represents text inputs.















