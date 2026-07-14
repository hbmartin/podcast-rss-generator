package transcript

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllErrorsAreConversionErrors(t *testing.T) {
	t.Parallel()
	errs := []error{
		newInvalidHTMLError(),
		newInvalidJSONError(),
		newInvalidXMLError(),
		newNoTranscriptFoundError(),
		newInvalidSRTError("bad block"),
		newInvalidVTTError(),
		newNoStartTimeError(),
		newUnknownFileTypeError("a/b.bin"),
	}
	for _, err := range errs {
		var conversionErr ConversionError
		require.ErrorAs(t, err, &conversionErr, "%T should be a ConversionError", err)
		assert.NotEmpty(t, err.Error())
	}
}

func TestInvalidXMLErrorMessage(t *testing.T) {
	t.Parallel()
	assert.Contains(t, newInvalidXMLError().Error(), "XML")
}

func TestInvalidJSONErrorMessage(t *testing.T) {
	t.Parallel()
	assert.Contains(t, newInvalidJSONError().Error(), "JSON")
}

func TestNoTranscriptFoundErrorMessage(t *testing.T) {
	t.Parallel()
	assert.Contains(t, newNoTranscriptFoundError().Error(), "transcript")
}

func TestInvalidSRTErrorKeepsBlock(t *testing.T) {
	t.Parallel()
	err := newInvalidSRTError("00:00 --> nonsense")
	assert.Equal(t, "00:00 --> nonsense", err.Block)
	assert.Contains(t, err.Error(), "00:00 --> nonsense")
}

func TestUnknownFileTypeErrorKeepsPath(t *testing.T) {
	t.Parallel()
	err := newUnknownFileTypeError("a/b.bin")
	assert.Equal(t, "a/b.bin", err.FilePath)
	assert.Contains(t, err.Error(), "a/b.bin")
}

func TestErrorsAsSpecificTypes(t *testing.T) {
	t.Parallel()
	var (
		htmlErr    *InvalidHTMLError
		srtErr     *InvalidSRTError
		unknownErr *UnknownFileTypeError
	)
	require.ErrorAs(t, error(newInvalidHTMLError()), &htmlErr)
	require.ErrorAs(t, error(newInvalidSRTError("x")), &srtErr)
	require.ErrorAs(t, error(newUnknownFileTypeError("p")), &unknownErr)
}
