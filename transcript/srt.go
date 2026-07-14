package transcript

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	srtSplitRe  = regexp.MustCompile(`\n\n\d*\n`)
	srtBlockRe  = regexp.MustCompile(`(?s)(\d+:\d+:\d+[,.]\d+)\s*-->\s*(\d+:\d+:\d+[,.]\d+)(\s*)(.*)`)
	mtsReplacer = strings.NewReplacer(",", ":", ".", ":")
)

// errInvalidTimestamp is the base error for malformed SRT/VTT timestamps.
var errInvalidTimestamp = errors.New("invalid timestamp")

const (
	// millisPerSecond scales SRT/VTT millisecond fields to fractional seconds.
	millisPerSecond = 1000
	// mtsParts is the number of colon-separated fields in HH:MM:SS:mmm.
	mtsParts = 4
)

// mtsToSecsFloat converts an HH:MM:SS[,.]mmm timestamp to seconds, rounded to
// three decimals.
func mtsToSecsFloat(timeString string) (float64, error) {
	parts := strings.Split(mtsReplacer.Replace(timeString), ":")
	if len(parts) != mtsParts {
		return 0, fmt.Errorf("%w: %q", errInvalidTimestamp, timeString)
	}
	values := make([]int, len(parts))
	for i, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return 0, fmt.Errorf("%w: %q", errInvalidTimestamp, timeString)
		}
		values[i] = value
	}
	total := float64(values[0]*3600+values[1]*60+values[2]) + float64(values[3])/millisPerSecond
	return math.Round(total*millisPerSecond) / millisPerSecond, nil
}

// ParseSRT parses an SRT transcript into a Transcript. It returns an
// *InvalidSRTError for a block that has no timing line.
func ParseSRT(srtString string) (*Transcript, error) {
	blocks := srtSplitRe.Split(srtString, -1)
	segments := make([]Segment, 0, len(blocks))
	for _, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		segment, err := srtBlockToSegment(block)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}
	return newTranscript(segments), nil
}

func srtBlockToSegment(block string) (Segment, error) {
	match := srtBlockRe.FindStringSubmatch(block)
	if match == nil {
		return Segment{}, newInvalidSRTError(block)
	}
	speaker, body := SplitSpeakerPrefix(strings.ReplaceAll(strings.TrimSpace(match[4]), "\n", " "))
	start, err := mtsToSecsFloat(match[1])
	if err != nil {
		return Segment{}, err
	}
	end, err := mtsToSecsFloat(match[2])
	if err != nil {
		return Segment{}, err
	}
	segment := Segment{Body: strPtr(body), StartTime: floatPtr(start), EndTime: floatPtr(end)}
	if speaker != "" {
		segment.Speaker = strPtr(speaker)
	}
	return segment, nil
}

// SRTFileToJSONFile converts an SRT file to PodcastIndex transcript JSON,
// merging optional metadata. Nothing is written when parsing fails.
func SRTFileToJSONFile(srtFile, jsonFile string, metadata map[string]string) error {
	srtString, err := ReadTextRobust(srtFile)
	if err != nil {
		return err
	}
	transcript, err := ParseSRT(srtString)
	if err != nil {
		return fmt.Errorf("%w: %s", err, srtFile)
	}
	if len(metadata) > 0 {
		transcript.Metadata = metadata
	}
	return writeTranscriptJSON(jsonFile, transcript)
}

// writeTranscriptJSON serializes a transcript to indented JSON and writes it.
func writeTranscriptJSON(jsonFile string, transcript *Transcript) error {
	data, err := encodeJSON(transcript, jsonIndent)
	if err != nil {
		return err
	}
	return WriteTextUTF8(jsonFile, string(data))
}

// jsonIndent is the indentation width used for transcript JSON files.
const jsonIndent = 4
