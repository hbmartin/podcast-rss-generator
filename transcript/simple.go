package transcript

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	secondsPerHour   = 3600
	secondsPerMinute = 60
)

var (
	errInvalidStartTime = errors.New("invalid startTime")
	errMissingSegments  = errors.New("missing segments")
	errInvalidSegment   = errors.New("invalid segment")
)

// numberToTs formats whole seconds as HH:MM:SS, truncating any fraction.
func numberToTs(seconds float64) string {
	total := int(seconds)
	hours := total / secondsPerHour
	minutes := (total % secondsPerHour) / secondsPerMinute
	secs := total % secondsPerMinute
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// segmentToLine renders a segment as "(HH:MM:SS) Speaker: body". It returns a
// *NoStartTimeError when the segment has no startTime.
func segmentToLine(segment map[string]any) (string, error) {
	startValue, ok := segment["startTime"]
	if !ok {
		return "", newNoStartTimeError()
	}
	seconds, err := toSeconds(startValue)
	if err != nil {
		return "", err
	}
	speaker := ""
	if value, ok := segment["speaker"]; ok {
		speaker = fmt.Sprintf("%v: ", value)
	}
	return fmt.Sprintf("(%s) %s%v", numberToTs(seconds), speaker, segment["body"]), nil
}

// toSeconds coerces a JSON startTime value to a float, matching Python's int().
func toSeconds(value any) (float64, error) {
	switch n := value.(type) {
	case float64:
		return n, nil
	case json.Number:
		return n.Float64()
	case string:
		return strconv.ParseFloat(n, 64)
	default:
		return 0, fmt.Errorf("%w: %v", errInvalidStartTime, value)
	}
}

// JSONFileToSimpleFile converts a PodcastIndex transcript JSON file to a simple
// one-line-per-segment text file. It returns a *NoStartTimeError when a segment
// lacks a startTime.
func JSONFileToSimpleFile(originFile, destinationFile string) error {
	text, err := ReadTextRobust(originFile)
	if err != nil {
		return err
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return fmt.Errorf("%s: %w", originFile, err)
	}
	rawSegments, ok := data["segments"].([]any)
	if !ok {
		return fmt.Errorf("%w: %s", errMissingSegments, originFile)
	}
	lines := make([]string, 0, len(rawSegments))
	for _, raw := range rawSegments {
		segment, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("%w: %s", errInvalidSegment, originFile)
		}
		line, err := segmentToLine(segment)
		if err != nil {
			return fmt.Errorf("%w: %s", err, originFile)
		}
		lines = append(lines, line)
	}
	return WriteTextUTF8(destinationFile, strings.Join(lines, "\n"))
}
