package chapters

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeID3StringEncodings(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "café", decodeID3String(encLatin1, []byte{0x63, 0x61, 0x66, 0xe9}))
	assert.Equal(t, "Hi", decodeID3String(encUTF8, []byte("Hi")))
	assert.Equal(t, "Hi", decodeID3String(encUTF16BE, []byte{0x00, 'H', 0x00, 'i'}))
	assert.Equal(t, "Hi", decodeID3String(encUTF16BOM, []byte{0xFF, 0xFE, 'H', 0x00, 'i', 0x00}))
	assert.Equal(t, "Hi", decodeID3String(encUTF16BOM, []byte{0xFE, 0xFF, 0x00, 'H', 0x00, 'i'}))
}

func TestDecodeTextFrameUTF16(t *testing.T) {
	t.Parallel()
	// encoding byte 2 (UTF-16BE) followed by "Hi".
	body := []byte{encUTF16BE, 0x00, 'H', 0x00, 'i'}
	assert.Equal(t, "Hi", decodeTextFrame(body))
}

func TestDecodeURLFrame(t *testing.T) {
	t.Parallel()
	// WOAR/WORS are plain URL link frames (Latin-1 URL, no encoding byte).
	assert.Equal(t, "https://example.com/a", decodeURLFrame("WOAR", []byte("https://example.com/a")))
	assert.Equal(t, "https://example.com/b", decodeURLFrame("WORS", []byte("https://example.com/b\x00")))
	// WXXX has an encoding byte + description + URL.
	wxxx := append([]byte{encUTF8}, []byte("desc\x00https://example.com/c")...)
	assert.Equal(t, "https://example.com/c", decodeURLFrame("WXXX", wxxx))
	assert.Empty(t, decodeURLFrame("WXXX", nil))
}

func TestCoerceInt(t *testing.T) {
	t.Parallel()
	assertOK := func(value any, want int) {
		got, ok := coerceInt(value)
		assert.True(t, ok)
		assert.Equal(t, want, got)
	}
	assertOK(float64(310), 310)
	assertOK(float64(3.9), 3)
	assertOK(310, 310)
	assertOK(int64(42), 42)
	assertOK("310", 310)
	assertOK(json.Number("7"), 7)

	for _, bad := range []any{"not-a-number", nil, true, []any{}} {
		_, ok := coerceInt(bad)
		assert.False(t, ok, "%v", bad)
	}
}
