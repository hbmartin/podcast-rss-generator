package chapters

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"sort"
	"strings"
	"unicode/utf16"
)

// ID3 parsing constants.
const (
	msPerSec = 1000

	id3HeaderSize   = 10
	frameHeaderSize = 10
	id3Version23    = 3
	id3Version24    = 4
	extHeaderFlag   = 0x40
	sizeFieldBytes  = 4
	synchsafeMask   = 0x7f
	utf16UnitBytes  = 2
	chapFieldsBytes = 16 // element ID terminator excluded; 4 x uint32 timing fields.

	// ID3 text encodings.
	encLatin1   = 0
	encUTF16BOM = 1
	encUTF16BE  = 2
	encUTF8     = 3
)

// errNoID3 signals that a file has no ID3v2 header (mirrors mutagen's
// ID3NoHeaderError, which the extractor maps to a nil result).
var errNoID3 = errors.New("no ID3 header")

// id3Frame is a decoded ID3v2 frame.
type id3Frame struct {
	id   string
	body []byte
}

// chapFrame is a decoded CHAP (chapter) frame.
type chapFrame struct {
	elementID string
	startTime uint32
	subFrames []id3Frame
}

// ctocFrame is a decoded CTOC (table of contents) frame.
type ctocFrame struct {
	childElementIDs []string
}

// ExtractID3Chapters extracts chapters from the ID3v2 CHAP frames of an MP3
// file. Chapter order follows the CTOC (table of contents) frame when present,
// otherwise chapters are sorted by start time. It returns nil when the file is
// missing, has no ID3 header, or contains no CHAP frames.
func ExtractID3Chapters(audioFile string) []Chapter {
	if !fileExists(audioFile) {
		return nil
	}
	data, err := os.ReadFile(audioFile) //nolint:gosec // reading a caller-supplied audio path is this function's purpose
	if err != nil {
		return nil
	}
	frames, err := parseID3(data)
	if err != nil {
		return nil
	}

	chapFrames := map[string]chapFrame{}
	var chapOrder []string
	var ctocs []ctocFrame
	for _, frame := range frames {
		switch frame.id {
		case "CHAP":
			chap, ok := parseCHAP(frame.body)
			if !ok {
				continue
			}
			if _, seen := chapFrames[chap.elementID]; !seen {
				chapOrder = append(chapOrder, chap.elementID)
			}
			chapFrames[chap.elementID] = chap
		case "CTOC":
			if ctoc, ok := parseCTOC(frame.body); ok {
				ctocs = append(ctocs, ctoc)
			}
		}
	}
	if len(chapFrames) == 0 {
		return nil
	}

	orderedIDs := ctocOrder(ctocs, chapFrames, chapOrder)
	chapters := make([]Chapter, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		chapters = append(chapters, chapToChapter(chapFrames[id]))
	}
	return chapters
}

// ctocOrder returns the chapter element IDs ordered per the first CTOC frame
// with known children, appending any unreferenced chapters sorted by start
// time. Without a usable CTOC, all chapters are sorted by start time.
func ctocOrder(ctocs []ctocFrame, chapFrames map[string]chapFrame, chapOrder []string) []string {
	sortedIDs := append([]string(nil), chapOrder...)
	sort.SliceStable(sortedIDs, func(i, j int) bool {
		return chapFrames[sortedIDs[i]].startTime < chapFrames[sortedIDs[j]].startTime
	})

	for _, ctoc := range ctocs {
		var childIDs []string
		for _, id := range ctoc.childElementIDs {
			if _, ok := chapFrames[id]; ok {
				childIDs = append(childIDs, id)
			}
		}
		if len(childIDs) == 0 {
			continue
		}
		childSet := make(map[string]struct{}, len(childIDs))
		for _, id := range childIDs {
			childSet[id] = struct{}{}
		}
		ordered := append([]string(nil), childIDs...)
		for _, id := range sortedIDs {
			if _, ok := childSet[id]; !ok {
				ordered = append(ordered, id)
			}
		}
		return ordered
	}
	return sortedIDs
}

// chapToChapter converts a CHAP frame to a Chapter, reading its title from the
// embedded TIT2 frame and its URL from an embedded WXXX/WOAR/WORS frame.
func chapToChapter(chap chapFrame) Chapter {
	title := ""
	url := ""
	for _, sub := range chap.subFrames {
		switch sub.id {
		case "TIT2":
			title = decodeTextFrame(sub.body)
		case "WXXX", "WOAR", "WORS":
			if u := decodeURLFrame(sub.id, sub.body); u != "" {
				url = u
			}
		}
	}
	if title == "" {
		title = chap.elementID
	}
	return Chapter{
		Start: int(chap.startTime) / msPerSec,
		Title: title,
		URL:   url,
	}
}

// parseID3 parses the ID3v2 header and returns the top-level frames. It supports
// ID3v2.3 (plain frame sizes) and ID3v2.4 (synchsafe frame sizes).
func parseID3(data []byte) ([]id3Frame, error) {
	if len(data) < id3HeaderSize || string(data[0:3]) != "ID3" {
		return nil, errNoID3
	}
	majorVersion := data[3]
	if majorVersion != id3Version23 && majorVersion != id3Version24 {
		return nil, errNoID3
	}
	flags := data[5]
	tagSize := int(synchsafe(data[6:id3HeaderSize]))
	end := min(id3HeaderSize+tagSize, len(data))
	body := data[id3HeaderSize:end]

	// Skip an extended header when present.
	if flags&extHeaderFlag != 0 && len(body) >= sizeFieldBytes {
		var extSize int
		if majorVersion == id3Version24 {
			extSize = int(synchsafe(body[0:sizeFieldBytes]))
		} else {
			extSize = int(binary.BigEndian.Uint32(body[0:sizeFieldBytes])) + sizeFieldBytes
		}
		if extSize <= len(body) {
			body = body[extSize:]
		}
	}

	return parseFrames(body, majorVersion == id3Version24), nil
}

// parseFrames parses a sequence of ID3v2.3/2.4 frames. When synchsafeSize is
// true, frame sizes are synchsafe (v2.4); otherwise plain big-endian (v2.3).
func parseFrames(buf []byte, synchsafeSize bool) []id3Frame {
	var frames []id3Frame
	pos := 0
	for pos+frameHeaderSize <= len(buf) {
		id := buf[pos : pos+sizeFieldBytes]
		if id[0] == 0 || !isFrameID(id) {
			break // padding or invalid
		}
		var size int
		if synchsafeSize {
			size = int(synchsafe(buf[pos+sizeFieldBytes : pos+2*sizeFieldBytes]))
		} else {
			size = int(binary.BigEndian.Uint32(buf[pos+sizeFieldBytes : pos+2*sizeFieldBytes]))
		}
		start := pos + frameHeaderSize
		if size < 0 || start+size > len(buf) {
			break
		}
		frames = append(frames, id3Frame{id: string(id), body: buf[start : start+size]})
		pos = start + size
	}
	return frames
}

// parseCHAP decodes a CHAP frame body.
func parseCHAP(body []byte) (chapFrame, bool) {
	elementID, rest, ok := readLatin1NulString(body)
	if !ok || len(rest) < chapFieldsBytes {
		return chapFrame{}, false
	}
	startTime := binary.BigEndian.Uint32(rest[0:sizeFieldBytes])
	subFrames := parseFrames(rest[chapFieldsBytes:], true)
	return chapFrame{elementID: elementID, startTime: startTime, subFrames: subFrames}, true
}

// parseCTOC decodes a CTOC frame body.
func parseCTOC(body []byte) (ctocFrame, bool) {
	_, rest, ok := readLatin1NulString(body)
	if !ok || len(rest) < utf16UnitBytes {
		return ctocFrame{}, false
	}
	entryCount := int(rest[1])
	remaining := rest[utf16UnitBytes:]
	childIDs := make([]string, 0, entryCount)
	for range entryCount {
		id, next, idOK := readLatin1NulString(remaining)
		if !idOK {
			break
		}
		childIDs = append(childIDs, id)
		remaining = next
	}
	return ctocFrame{childElementIDs: childIDs}, true
}

// decodeTextFrame decodes a text frame (for example TIT2), returning its first
// text value.
func decodeTextFrame(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	text := decodeID3String(body[0], body[1:])
	if i := strings.IndexByte(text, 0); i >= 0 {
		text = text[:i]
	}
	return text
}

// decodeURLFrame decodes the URL from a WXXX (encoding + description + URL) or a
// plain URL link frame (WOAR/WORS).
func decodeURLFrame(id string, body []byte) string {
	if id == "WXXX" {
		if len(body) == 0 {
			return ""
		}
		encoding := body[0]
		_, rest := splitID3NulString(encoding, body[1:])
		return strings.Trim(decodeLatin1(rest), "\x00")
	}
	return strings.Trim(decodeLatin1(body), "\x00")
}

// synchsafe decodes a 4-byte synchsafe integer (7 bits per byte).
func synchsafe(b []byte) uint32 {
	return uint32(b[0]&synchsafeMask)<<21 |
		uint32(b[1]&synchsafeMask)<<14 |
		uint32(b[2]&synchsafeMask)<<7 |
		uint32(b[3]&synchsafeMask)
}

// isFrameID reports whether id is a plausible ID3v2 frame identifier
// (uppercase letters and digits).
func isFrameID(id []byte) bool {
	for _, c := range id {
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// readLatin1NulString reads a NUL-terminated Latin-1 string, returning the
// string, the remaining bytes, and whether a terminator was found.
func readLatin1NulString(data []byte) (string, []byte, bool) {
	before, after, found := bytes.Cut(data, []byte{0})
	if !found {
		return "", nil, false
	}
	return decodeLatin1(before), after, true
}

// splitID3NulString splits data at the NUL terminator appropriate for the
// encoding, returning the (decoded) first segment and the remaining bytes.
func splitID3NulString(encoding byte, data []byte) (string, []byte) {
	if encoding == encUTF16BOM || encoding == encUTF16BE {
		for i := 0; i+1 < len(data); i += utf16UnitBytes {
			if data[i] == 0 && data[i+1] == 0 {
				return decodeID3String(encoding, data[:i]), data[i+utf16UnitBytes:]
			}
		}
		return decodeID3String(encoding, data), nil
	}
	if before, after, found := bytes.Cut(data, []byte{0}); found {
		return decodeID3String(encoding, before), after
	}
	return decodeID3String(encoding, data), nil
}

// decodeID3String decodes ID3 text bytes per the encoding byte: 0=Latin-1,
// 1=UTF-16 with BOM, 2=UTF-16BE, 3=UTF-8.
func decodeID3String(encoding byte, data []byte) string {
	switch encoding {
	case encUTF16BOM:
		return decodeUTF16(data, true)
	case encUTF16BE:
		return decodeUTF16(data, false)
	case encUTF8:
		return string(data)
	case encLatin1:
		return decodeLatin1(data)
	default:
		return decodeLatin1(data)
	}
}

// decodeLatin1 decodes ISO-8859-1 bytes to a string.
func decodeLatin1(data []byte) string {
	runes := make([]rune, len(data))
	for i, b := range data {
		runes[i] = rune(b)
	}
	return string(runes)
}

// decodeUTF16 decodes UTF-16 bytes. When hasBOM is true a leading byte-order
// mark selects endianness (defaulting to little-endian); otherwise big-endian
// is assumed.
func decodeUTF16(data []byte, hasBOM bool) string {
	if len(data) < utf16UnitBytes {
		return ""
	}
	bigEndian := !hasBOM
	if hasBOM {
		switch {
		case data[0] == 0xFE && data[1] == 0xFF:
			bigEndian = true
			data = data[utf16UnitBytes:]
		case data[0] == 0xFF && data[1] == 0xFE:
			bigEndian = false
			data = data[utf16UnitBytes:]
		}
	}
	units := make([]uint16, 0, len(data)/utf16UnitBytes)
	for i := 0; i+1 < len(data); i += utf16UnitBytes {
		if bigEndian {
			units = append(units, uint16(data[i])<<8|uint16(data[i+1]))
		} else {
			units = append(units, uint16(data[i+1])<<8|uint16(data[i]))
		}
	}
	return string(utf16.Decode(units))
}
