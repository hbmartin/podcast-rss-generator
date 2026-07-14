package transcript

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"unicode/utf8"
)

// utf8BOM is the UTF-8 byte-order mark stripped from robustly read text.
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// outputFilePerm is the permission used for written transcript files.
const outputFilePerm = 0o644

// maxReadParallelism bounds concurrent file reads in MapFilesInParallel.
const maxReadParallelism = 16

// isFileAllowed reports whether a filename should be considered for conversion:
// it must not be in the ignore list, be hidden, or be a PDF or octet-stream.
func isFileAllowed(filename string, ignore []string) bool {
	if slices.Contains(ignore, filename) {
		return false
	}
	return !strings.HasPrefix(filename, ".") &&
		!strings.HasSuffix(filename, ".pdf") &&
		!strings.HasSuffix(filename, ".octet-stream")
}

// ListFiles walks directory recursively and returns the paths of files that
// pass isFileAllowed.
func ListFiles(directory string, ignore []string) ([]string, error) {
	var filePaths []string
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if isFileAllowed(d.Name(), ignore) {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return filePaths, nil
}

// ReadTextRobust reads a text file as UTF-8, stripping a leading BOM and
// replacing undecodable bytes with the Unicode replacement character, so
// real-world transcript files that are not valid UTF-8 do not abort a run.
func ReadTextRobust(filePath string) (string, error) {
	data, err := os.ReadFile(filePath) //nolint:gosec // reading a caller-supplied transcript path is this function's purpose
	if err != nil {
		return "", err
	}
	return decodeTextRobust(data), nil
}

// ReadFirstLine reads the first line (including its trailing newline, if any) of
// a file, applying the same BOM stripping and replacement as ReadTextRobust.
func ReadFirstLine(filePath string) (string, error) {
	data, err := os.ReadFile(filePath) //nolint:gosec // reading a caller-supplied transcript path is this function's purpose
	if err != nil {
		return "", err
	}
	text := decodeTextRobust(data)
	if i := strings.IndexByte(text, '\n'); i >= 0 {
		return text[:i+1], nil
	}
	return text, nil
}

// WriteTextUTF8 writes data to filePath as UTF-8.
func WriteTextUTF8(filePath, data string) error {
	return os.WriteFile(filePath, []byte(data), outputFilePerm)
}

// decodeTextRobust strips a leading UTF-8 BOM, replaces each invalid byte with
// the Unicode replacement character, and applies universal-newline translation
// (\r\n and \r become \n), matching Python's text-mode read_text.
func decodeTextRobust(data []byte) string {
	data = stripBOM(data)
	var text string
	if utf8.Valid(data) {
		text = string(data)
	} else {
		var b strings.Builder
		b.Grow(len(data))
		for i := 0; i < len(data); {
			r, size := utf8.DecodeRune(data[i:])
			if r == utf8.RuneError && size == 1 {
				b.WriteRune(utf8.RuneError)
				i++
				continue
			}
			b.WriteRune(r)
			i += size
		}
		text = b.String()
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.ReplaceAll(text, "\r", "\n")
}

func stripBOM(data []byte) []byte {
	if len(data) >= len(utf8BOM) &&
		data[0] == utf8BOM[0] && data[1] == utf8BOM[1] && data[2] == utf8BOM[2] {
		return data[len(utf8BOM):]
	}
	return data
}

// FileResult pairs an input path with the result of applying a transform to it.
type FileResult[T any] struct {
	Path   string
	Result T
}

// MapFilesInParallel applies transform to each path concurrently (bounded), and
// returns the results in input order.
func MapFilesInParallel[T any](
	ctx context.Context,
	paths []string,
	transform func(context.Context, string) T,
) []FileResult[T] {
	if len(paths) == 0 {
		return nil
	}
	results := make([]FileResult[T], len(paths))
	limit := min(len(paths), maxReadParallelism)
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	for i := range paths {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = FileResult[T]{Path: paths[i], Result: transform(ctx, paths[i])}
		}()
	}
	wg.Wait()
	return results
}
