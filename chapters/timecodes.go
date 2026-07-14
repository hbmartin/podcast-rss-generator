package chapters

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Conversions between HH:MM:SS-style timestamps and whole seconds.
const (
	secsPerMinute = 60
	secsPerHour   = 3600
	maxTimeParts  = 3
	hoursParts    = 3
	minutesParts  = 2
)

// errInvalidTimestamp is the base error for malformed timestamps; errNegativeSeconds
// for negative second counts passed to SecsToTs.
var (
	errInvalidTimestamp = errors.New("invalid timestamp")
	errNegativeSeconds  = errors.New("seconds must be non-negative")
)

// TsToSecs converts a timestamp string ([[HH:]MM:]SS[.fff]) to whole seconds.
// Fractional seconds are truncated. It returns an error for strings that are
// not valid timestamps (non-numeric parts, too many segments, or
// minutes/seconds >= 60 when a larger unit is present).
func TsToSecs(timeString string) (int, error) {
	rawTime := strings.TrimSpace(timeString)
	wholeSeconds, fractionalSeconds, separator := strings.Cut(rawTime, ".")
	if separator {
		if strings.Contains(fractionalSeconds, ":") {
			return 0, fmt.Errorf("%w: fractional seconds must be in final segment: %q", errInvalidTimestamp, timeString)
		}
		if !isDigits(fractionalSeconds) {
			return 0, fmt.Errorf("%w: non-numeric fractional seconds: %q", errInvalidTimestamp, timeString)
		}
	}

	timeParts := strings.Split(wholeSeconds, ":")
	if len(timeParts) > maxTimeParts {
		return 0, fmt.Errorf("%w: too many segments: %q", errInvalidTimestamp, timeString)
	}

	values := make([]int, len(timeParts))
	for i, part := range timeParts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return 0, fmt.Errorf("%w: non-numeric segment: %q", errInvalidTimestamp, timeString)
		}
		values[i] = value
	}

	for _, value := range values {
		if value < 0 {
			return 0, fmt.Errorf("%w: negative segment: %q", errInvalidTimestamp, timeString)
		}
	}
	for _, value := range values[1:] {
		if value >= secsPerMinute {
			return 0, fmt.Errorf("%w: minutes and seconds must be < 60: %q", errInvalidTimestamp, timeString)
		}
	}

	seconds := values[len(values)-1]
	if len(values) >= minutesParts {
		seconds += secsPerMinute * values[len(values)-minutesParts]
	}
	if len(values) >= hoursParts {
		seconds += secsPerHour * values[len(values)-hoursParts]
	}
	return seconds, nil
}

// SecsToTs converts whole seconds to a timestamp (M:SS, or H:MM:SS at or above
// one hour). It returns an error for negative input.
func SecsToTs(seconds int) (string, error) {
	if seconds < 0 {
		return "", fmt.Errorf("%w: %d", errNegativeSeconds, seconds)
	}
	hours := seconds / secsPerHour
	remainder := seconds % secsPerHour
	minutes := remainder / secsPerMinute
	secs := remainder % secsPerMinute
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs), nil
	}
	return fmt.Sprintf("%d:%02d", minutes, secs), nil
}

// isDigits reports whether s is non-empty and all ASCII digits, matching
// Python's str.isdigit for the inputs seen here.
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
