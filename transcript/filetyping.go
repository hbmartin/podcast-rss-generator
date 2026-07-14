package transcript

import (
	"context"
	"strings"
)

// FileType identifies a recognized transcript file format.
type FileType int

// FileType values.
const (
	FileTypeHTML FileType = iota
	FileTypeJSON
	FileTypeSRT
	FileTypeVTT
	FileTypeXML
	FileTypeUnknown
)

// String returns the lowercase name of the FileType.
func (t FileType) String() string {
	switch t {
	case FileTypeHTML:
		return "html"
	case FileTypeJSON:
		return "json"
	case FileTypeSRT:
		return "srt"
	case FileTypeVTT:
		return "vtt"
	case FileTypeXML:
		return "xml"
	case FileTypeUnknown:
		return "unknown"
	}
	return "unknown"
}

// allFileTypes lists every FileType, used to seed grouping maps.
var allFileTypes = []FileType{
	FileTypeHTML, FileTypeJSON, FileTypeSRT, FileTypeVTT, FileTypeXML, FileTypeUnknown,
}

// extensionToType maps a lowercase file extension to its FileType.
var extensionToType = map[string]FileType{
	"vtt":  FileTypeVTT,
	"srt":  FileTypeSRT,
	"htm":  FileTypeHTML,
	"html": FileTypeHTML,
	"json": FileTypeJSON,
	"xml":  FileTypeXML,
	"xsl":  FileTypeXML,
}

// extractFileTypeFromName determines a FileType from a path's extension.
func extractFileTypeFromName(filePath string) FileType {
	extension := filePath
	if i := strings.LastIndex(filePath, "."); i >= 0 {
		extension = filePath[i+1:]
	}
	if fileType, ok := extensionToType[strings.ToLower(extension)]; ok {
		return fileType
	}
	return FileTypeUnknown
}

// extractFileTypeFromFirstLine determines a FileType by inspecting the first
// line of a file.
func extractFileTypeFromFirstLine(line string) FileType {
	stripped := strings.TrimSpace(line)
	if stripped == "" {
		return FileTypeUnknown
	}
	if strings.Contains(line, "WEBVTT") {
		return FileTypeVTT
	}
	if stripped[0] == '1' || strings.Contains(line, "-->") {
		return FileTypeSRT
	}
	if strings.HasPrefix(stripped, "<?xml") || strings.Contains(line, "podlove.org/simple-transcripts") {
		return FileTypeXML
	}
	if strings.Contains(line, "<") {
		return FileTypeHTML
	}
	if strings.Contains(line, "{") && !strings.Contains(line, "rtf") {
		return FileTypeJSON
	}
	return FileTypeUnknown
}

// extractFileTypeFromContent determines a FileType from a file's first line,
// returning FileTypeUnknown when the file cannot be read.
func extractFileTypeFromContent(_ context.Context, filePath string) FileType {
	firstLine, err := ReadFirstLine(filePath)
	if err != nil {
		return FileTypeUnknown
	}
	return extractFileTypeFromFirstLine(firstLine)
}

// IdentifyFileType determines a file's transcript FileType from its extension,
// falling back to inspecting its first line.
func IdentifyFileType(ctx context.Context, filePath string) FileType {
	if fileType := extractFileTypeFromName(filePath); fileType != FileTypeUnknown {
		return fileType
	}
	return extractFileTypeFromContent(ctx, filePath)
}

// IdentifyFileTypes groups file paths by detected transcript type. Files whose
// extension is unrecognized have their first line inspected in parallel. The
// returned map has an entry for every FileType; unidentifiable files are under
// FileTypeUnknown.
func IdentifyFileTypes(ctx context.Context, filePaths []string) map[FileType][]string {
	grouped := make(map[FileType][]string, len(allFileTypes))
	for _, fileType := range allFileTypes {
		grouped[fileType] = []string{}
	}
	for _, filePath := range filePaths {
		fileType := extractFileTypeFromName(filePath)
		grouped[fileType] = append(grouped[fileType], filePath)
	}

	typesFromContent := MapFilesInParallel(ctx, grouped[FileTypeUnknown], extractFileTypeFromContent)
	grouped[FileTypeUnknown] = []string{}
	for _, result := range typesFromContent {
		grouped[result.Result] = append(grouped[result.Result], result.Path)
	}
	return grouped
}
