package transcript

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// VTT parsing regexes, mirroring webvtt-py.
var (
	vttCueTimingsRe = regexp.MustCompile(`^\s*((?:\d+:)?\d{2}:\d{2}.\d{3})\s*-->\s*((?:\d+:)?\d{2}:\d{2}.\d{3})`)
	vttTimestampRe  = regexp.MustCompile(`^(?:(\d+):)?(\d{1,2}):(\d{1,2})[.,](\d{3})`)
	vttVoiceRe      = regexp.MustCompile(`^<v(?:\.\w+)*\s+([^>]+)>`)
	vttTagRe        = regexp.MustCompile(`<.*?>`)
)

// vttFields is the number of captured components in a VTT timestamp
// (hours, minutes, seconds, milliseconds).
const vttFields = 4

// ParseVTT parses a WebVTT transcript into a Transcript. It returns an
// *InvalidVTTError when the content is not valid WebVTT.
func ParseVTT(vttString string) (*Transcript, error) {
	lines := splitLines(vttString)
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "WEBVTT") {
		return nil, newInvalidVTTError()
	}
	var segments []Segment
	for _, block := range iterBlocks(lines) {
		if !isValidCue(block) {
			continue
		}
		segment, err := cueToSegment(block)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}
	return newTranscript(segments), nil
}

// splitLines splits text into lines the way Python's str.splitlines does: on
// \n, \r\n, and \r, without a trailing empty element for a final line break.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// iterBlocks groups consecutive non-blank lines into blocks separated by blank
// lines, mirroring webvtt-py's iter_blocks_of_lines.
func iterBlocks(lines []string) [][]string {
	var blocks [][]string
	var current []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			current = append(current, line)
		} else if len(current) > 0 {
			blocks = append(blocks, current)
			current = nil
		}
	}
	if len(current) > 0 {
		blocks = append(blocks, current)
	}
	return blocks
}

// isValidCue reports whether a block of lines is a valid WebVTT cue block.
func isValidCue(lines []string) bool {
	if len(lines) >= 2 && vttCueTimingsRe.MatchString(lines[0]) && !strings.Contains(lines[1], "-->") {
		return true
	}
	return len(lines) >= 3 && !strings.Contains(lines[0], "-->") &&
		vttCueTimingsRe.MatchString(lines[1]) && !strings.Contains(lines[2], "-->")
}

// cueToSegment parses a valid cue block into a Segment.
func cueToSegment(lines []string) (Segment, error) {
	var startRaw, endRaw string
	var payload []string
	for _, line := range lines {
		if match := vttCueTimingsRe.FindStringSubmatch(line); match != nil {
			startRaw, endRaw = match[1], match[2]
		} else if startRaw == "" {
			continue // cue identifier line
		} else {
			payload = append(payload, line)
		}
	}

	text := vttTagRe.ReplaceAllString(strings.Join(payload, "\n"), "")
	body := strings.ReplaceAll(strings.TrimSpace(text), "\n", " ")
	speaker := cueVoice(payload)
	if speaker == "" {
		speaker, body = SplitSpeakerPrefix(body)
	}

	start, err := vttTimestampToSecs(startRaw)
	if err != nil {
		return Segment{}, err
	}
	end, err := vttTimestampToSecs(endRaw)
	if err != nil {
		return Segment{}, err
	}
	segment := Segment{Body: strPtr(body), StartTime: floatPtr(start), EndTime: floatPtr(end)}
	if speaker != "" {
		segment.Speaker = strPtr(speaker)
	}
	return segment, nil
}

// cueVoice returns the voice-span name from a cue's first payload line, or "".
func cueVoice(payload []string) string {
	if len(payload) > 0 && strings.HasPrefix(payload[0], "<v") {
		if match := vttVoiceRe.FindStringSubmatch(payload[0]); match != nil {
			return match[1]
		}
	}
	return ""
}

// vttTimestampToSecs converts a WebVTT cue timestamp ([HH:]MM:SS.mmm) to
// seconds, rounded to three decimals.
func vttTimestampToSecs(ts string) (float64, error) {
	match := vttTimestampRe.FindStringSubmatch(ts)
	if match == nil {
		return 0, fmt.Errorf("%w: %q", errInvalidTimestamp, ts)
	}
	var parts [vttFields]int
	for i, group := range match[1:] {
		if group == "" {
			continue
		}
		value, err := strconv.Atoi(group)
		if err != nil {
			return 0, fmt.Errorf("%w: %q: %w", errInvalidTimestamp, ts, err)
		}
		parts[i] = value
	}
	total := float64(parts[0]*3600+parts[1]*60+parts[2]) + float64(parts[3])/millisPerSecond
	return math.Round(total*millisPerSecond) / millisPerSecond, nil
}

// VTTFileToJSONFile converts a WebVTT file to PodcastIndex transcript JSON,
// merging optional metadata. Nothing is written when parsing fails.
func VTTFileToJSONFile(vttFile, jsonFile string, metadata map[string]string) error {
	vttString, err := ReadTextRobust(vttFile)
	if err != nil {
		return err
	}
	transcript, err := ParseVTT(vttString)
	if err != nil {
		return fmt.Errorf("%w: %s", err, vttFile)
	}
	if len(metadata) > 0 {
		transcript.Metadata = metadata
	}
	return writeTranscriptJSON(jsonFile, transcript)
}
