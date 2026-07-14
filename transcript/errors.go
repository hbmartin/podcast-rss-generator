package transcript

// ConversionError is implemented by every error raised while converting a
// transcript. Callers can test for it with errors.As and a
// ConversionError target.
type ConversionError interface {
	error
	isConversionError()
}

// baseConversionError carries the message shared by all conversion errors.
type baseConversionError struct {
	msg string
}

func (e baseConversionError) Error() string      { return e.msg }
func (e baseConversionError) isConversionError() {}

// InvalidHTMLError reports that HTML does not conform to the expected transcript
// format (no <cite> or <time> tags).
type InvalidHTMLError struct{ baseConversionError }

// InvalidJSONError reports that JSON cannot be parsed as a PodcastIndex
// transcript.
type InvalidJSONError struct{ baseConversionError }

// InvalidXMLError reports that XML does not conform to
// http://podlove.org/simple-transcripts.
type InvalidXMLError struct{ baseConversionError }

// NoTranscriptFoundError reports that no transcript blocks could be located in
// the source (for example an empty or unrecognized file).
type NoTranscriptFoundError struct{ baseConversionError }

// InvalidVTTError reports that WebVTT could not be parsed.
type InvalidVTTError struct{ baseConversionError }

// NoStartTimeError reports a missing startTime in the source transcript.
type NoStartTimeError struct{ baseConversionError }

// InvalidSRTError reports that an SRT block could not be parsed. Block holds the
// offending block.
type InvalidSRTError struct {
	baseConversionError
	Block string
}

// UnknownFileTypeError reports that a file's transcript format could not be
// determined. FilePath holds the offending path.
type UnknownFileTypeError struct {
	baseConversionError
	FilePath string
}

func newInvalidHTMLError() *InvalidHTMLError {
	return &InvalidHTMLError{baseConversionError{"the provided HTML file is not a valid transcript"}}
}

func newInvalidJSONError() *InvalidJSONError {
	return &InvalidJSONError{baseConversionError{"the provided JSON file is not a valid transcript"}}
}

func newInvalidXMLError() *InvalidXMLError {
	return &InvalidXMLError{baseConversionError{"the provided XML file is not a valid transcript"}}
}

func newNoTranscriptFoundError() *NoTranscriptFoundError {
	return &NoTranscriptFoundError{
		baseConversionError{"the provided source does not contain a transcript or could not be parsed"},
	}
}

func newInvalidVTTError() *InvalidVTTError {
	return &InvalidVTTError{baseConversionError{"the provided VTT could not be parsed"}}
}

func newNoStartTimeError() *NoStartTimeError {
	return &NoStartTimeError{baseConversionError{"failed to find startTime in source transcript"}}
}

func newInvalidSRTError(block string) *InvalidSRTError {
	return &InvalidSRTError{
		baseConversionError: baseConversionError{"the provided SRT could not be parsed:\n" + block},
		Block:               block,
	}
}

func newUnknownFileTypeError(filePath string) *UnknownFileTypeError {
	return &UnknownFileTypeError{
		baseConversionError: baseConversionError{"failed to determine the file type of " + filePath},
		FilePath:            filePath,
	}
}
