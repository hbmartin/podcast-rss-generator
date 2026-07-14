package chapters_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/hbmartin/podcast-rss-generator/v2/chapters"
	"github.com/stretchr/testify/assert"
)

func TestWithHeadersAndTimeout(t *testing.T) {
	t.Parallel()
	var gotHeader string
	var hadDeadline bool
	client := &fakeClient{do: func(req *http.Request) (*http.Response, error) {
		gotHeader = req.Header.Get("User-Agent")
		_, hadDeadline = req.Context().Deadline()
		return newResponse(http.StatusOK, mustMarshal(pciJSONDoc(t))), nil
	}}

	got := chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json",
		chapters.WithHeaders(map[string]string{"User-Agent": "podcast-tools"}),
		chapters.WithTimeout(5*time.Second),
	)
	assert.Equal(t, expectedPCI(), got)
	assert.Equal(t, "podcast-tools", gotHeader)
	assert.True(t, hadDeadline)
}

func TestWithoutTimeoutHasNoDeadline(t *testing.T) {
	t.Parallel()
	var hadDeadline bool
	client := &fakeClient{do: func(req *http.Request) (*http.Response, error) {
		_, hadDeadline = req.Context().Deadline()
		return newResponse(http.StatusOK, mustMarshal(pciJSONDoc(t))), nil
	}}
	chapters.GetAndExtractPCIChapters(
		context.Background(), client, "https://example.com/chapters.json",
		chapters.WithTimeout(0),
	)
	assert.False(t, hadDeadline)
}
