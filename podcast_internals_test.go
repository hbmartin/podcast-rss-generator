package podcast

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

type valueWriter struct{}

func (valueWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestStringError(t *testing.T) {
	t.Parallel()

	// arrange
	p := Podcast{}
	p.encode = func(_ io.Writer, _ any) error {
		return io.ErrClosedPipe
	}

	// act
	r := p.String()

	// assert
	assert.Contains(t, r, io.ErrClosedPipe.Error())
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

func TestIsNilWriter(t *testing.T) {
	t.Parallel()

	var nilBuffer *bytes.Buffer

	tests := []struct {
		name string
		w    io.Writer
		want bool
	}{
		{
			name: "nil",
			w:    nil,
			want: true,
		},
		{
			name: "typed nil",
			w:    nilBuffer,
			want: true,
		},
		{
			name: "pointer writer",
			w:    &bytes.Buffer{},
		},
		{
			name: "value writer",
			w:    valueWriter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isNilWriter(tt.w))
		})
	}
}

func TestTruncateRunes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		in    string
		limit int
		want  string
	}{
		{
			name:  "non-positive limit",
			in:    "abc",
			limit: 0,
			want:  "",
		},
		{
			name:  "within limit",
			in:    "épisode",
			limit: 7,
			want:  "épisode",
		},
		{
			name:  "truncates by runes",
			in:    "épisode",
			limit: 3,
			want:  "épi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, truncateRunes(tt.in, tt.limit))
		})
	}
}

func TestTruncateRunesWithSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     string
		limit  int
		suffix string
		want   string
	}{
		{
			name:   "within limit",
			in:     "short",
			limit:  10,
			suffix: "...",
			want:   "short",
		},
		{
			name:   "suffix reaches limit",
			in:     "episode",
			limit:  2,
			suffix: "...",
			want:   "..",
		},
		{
			name:   "truncates with suffix",
			in:     "episode",
			limit:  5,
			suffix: "...",
			want:   "ep...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, truncateRunesWithSuffix(tt.in, tt.limit, tt.suffix))
		})
	}
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

func TestClonePointerSlice(t *testing.T) {
	t.Parallel()

	// empty and nil slices are returned as-is.
	assert.Nil(t, clonePointerSlice[PPerson](nil))
	assert.Empty(t, clonePointerSlice([]*PPerson{}))

	// entries are deep-copied and nil entries are dropped.
	original := &PPerson{Name: "Jane Doe"}
	cloned := clonePointerSlice([]*PPerson{original, nil})

	assert.Len(t, cloned, 1)
	assert.NotSame(t, original, cloned[0])

	original.Name = "mutated"
	assert.Equal(t, "Jane Doe", cloned[0].Name)
}

func TestHasPodcastElementsSkipsNilItems(t *testing.T) {
	t.Parallel()

	p := New("t", "l", "d", time.Time{}, time.Time{})
	p.Items = append(p.Items, nil)
	assert.False(t, p.hasPodcastElements())

	p.Items = append(p.Items, &Item{PChapters: &PChapters{URL: "u", Type: "t"}})
	assert.True(t, p.hasPodcastElements())
}
