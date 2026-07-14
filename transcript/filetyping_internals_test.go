package transcript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFileTypeFromName(t *testing.T) {
	t.Parallel()
	cases := map[string]FileType{
		"a/b.vtt":  FileTypeVTT,
		"a/b.srt":  FileTypeSRT,
		"a/b.HTML": FileTypeHTML,
		"a/b.htm":  FileTypeHTML,
		"a/b.json": FileTypeJSON,
		"a/b.xml":  FileTypeXML,
		"a/b.xsl":  FileTypeXML,
		"a/b.txt":  FileTypeUnknown,
		"a/b":      FileTypeUnknown,
	}
	for input, want := range cases {
		assert.Equal(t, want, extractFileTypeFromName(input), input)
	}
}

func TestExtractFileTypeFromFirstLine(t *testing.T) {
	t.Parallel()
	cases := map[string]FileType{
		"WEBVTT\n":                        FileTypeVTT,
		"1\n":                             FileTypeSRT,
		"00:00:00,000 --> 00:00:01,000\n": FileTypeSRT,
		"<!DOCTYPE html>\n":               FileTypeHTML,
		`<?xml version="1.0"?>`:           FileTypeXML,
		`{"version": "1.0.0"}`:            FileTypeJSON,
		"hello there":                     FileTypeUnknown,
	}
	for input, want := range cases {
		assert.Equal(t, want, extractFileTypeFromFirstLine(input), input)
	}
}

func TestExtractFileTypeFromFirstLineEmpty(t *testing.T) {
	t.Parallel()
	assert.Equal(t, FileTypeUnknown, extractFileTypeFromFirstLine(""))
	assert.Equal(t, FileTypeUnknown, extractFileTypeFromFirstLine("\n"))
	assert.Equal(t, FileTypeUnknown, extractFileTypeFromFirstLine("   \n"))
}
