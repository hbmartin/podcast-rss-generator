package podcast

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestStringError(t *testing.T) {
	t.Parallel()

	// arrange
	e := "TestEncodeError error result"
	p := Podcast{}
	p.encode = func(w io.Writer, o interface{}) error {
		return errors.New(e)
	}

	// act
	r := p.String()

	// assert
	assert.Contains(t, r, e)
}

func TestEncodeError(t *testing.T) {
	t.Parallel()

	// arrange
	p := New("", "", "", time.Time{}, time.Time{})
	b := []byte{}
	w := bytes.NewBuffer(b)
	c := new(chan bool)

	// act
	err := p.encode(w, c)

	// assert
	assert.Error(t, err)
}

func TestParseDuration(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "0:00", parseDuration(0))
	assert.Equal(t, "0:40", parseDuration(40))
	assert.Equal(t, "1:00", parseDuration(60))
	assert.Equal(t, "1:40", parseDuration(100))
	assert.Equal(t, "2:01", parseDuration(121))
	assert.Equal(t, "59:59", parseDuration(3599))
	assert.Equal(t, "1:00:00", parseDuration(3600))
	assert.Equal(t, "1:00:01", parseDuration(3601))
	assert.Equal(t, "1:01:00", parseDuration(3660))
	assert.Equal(t, "1:01:03", parseDuration(3663))
	assert.Equal(t, "10:00:00", parseDuration(36000))
	assert.Equal(t, "10:00:01", parseDuration(36001))
	assert.Equal(t, "10:01:00", parseDuration(36060))
	assert.Equal(t, "10:01:03", parseDuration(36063))
}

func TestParseDescriptionByteLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "keeps text under limit",
			in:   "short description",
			want: "short description",
		},
		{
			name: "truncates ascii at byte limit",
			in:   strings.Repeat("a", 4001),
			want: strings.Repeat("a", 4000),
		},
		{
			name: "keeps multibyte rune ending at byte limit",
			in:   strings.Repeat("a", 3998) + "é" + "b",
			want: strings.Repeat("a", 3998) + "é",
		},
		{
			name: "drops multibyte rune crossing byte limit",
			in:   strings.Repeat("a", 3999) + "é",
			want: strings.Repeat("a", 3999),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := string(parseDescription(tt.in))
			if got != tt.want {
				t.Fatalf("parseDescription() = %q, want %q", got, tt.want)
			}
			if len(got) > 4000 {
				t.Fatalf("parseDescription() byte length = %d, want <= 4000", len(got))
			}
			if !utf8.ValidString(got) {
				t.Fatal("parseDescription() returned invalid UTF-8")
			}
		})
	}
}
