package transcript

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// dirPerm is the permission for directories created during bulk conversion.
const dirPerm = 0o750

// errNoDBLister reports a ".db" source path with no DBLister configured.
var errNoDBLister = errors.New("no DB lister configured")

// converterFunc converts a source transcript file to a destination JSON file,
// merging optional metadata.
type converterFunc func(origin, destination string, metadata map[string]string) error

// converterOrder lists the FileTypes with converters, in a stable order.
var converterOrder = []FileType{
	FileTypeHTML, FileTypeJSON, FileTypeSRT, FileTypeVTT, FileTypeXML,
}

// fileConverters returns the converter for each supported FileType.
func fileConverters() map[FileType]converterFunc {
	return map[FileType]converterFunc{
		FileTypeHTML: HTMLFileToJSONFile,
		FileTypeJSON: JSONFileToJSONFile,
		FileTypeSRT:  SRTFileToJSONFile,
		FileTypeVTT:  VTTFileToJSONFile,
		FileTypeXML:  XMLFileToJSONFile,
	}
}

// ConversionPair is a (source, destination) pair.
type ConversionPair struct {
	Source      string
	Destination string
}

// ConversionFailure is a (source, error message) pair.
type ConversionFailure struct {
	Source  string
	Message string
}

// ConversionSummary is the outcome of a BulkConvert run.
//
//   - Converted: pairs successfully written (or planned, when DryRun is true).
//   - Skipped: pairs whose destination already existed and Overwrite was false.
//   - Failed: sources that could not be converted, with the error message.
//   - Unknown: sources whose transcript format could not be identified.
type ConversionSummary struct {
	Converted []ConversionPair
	Skipped   []ConversionPair
	Failed    []ConversionFailure
	Unknown   []string
	DryRun    bool
}

// DBLister lists transcript file paths and per-file metadata from a database
// (for example an overcast-to-sqlite database). It is injected so the core
// package carries no database dependency.
type DBLister interface {
	ListFiles(ctx context.Context, dbPath string, ignore []string) (paths []string, metadata map[string]map[string]string, err error)
}

// bulkConfig holds resolved BulkConvert options.
type bulkConfig struct {
	ignore     []string
	overwrite  bool
	dryRun     bool
	dbLister   DBLister
	fileLister func(directory string, ignore []string) ([]string, error)
}

// BulkOption configures BulkConvert.
type BulkOption func(*bulkConfig)

// WithIgnore skips files whose base name appears in ignore.
func WithIgnore(ignore []string) BulkOption {
	return func(c *bulkConfig) { c.ignore = ignore }
}

// WithOverwrite re-converts files whose destination JSON already exists.
func WithOverwrite() BulkOption {
	return func(c *bulkConfig) { c.overwrite = true }
}

// WithDryRun reports what would be converted without writing any files.
func WithDryRun() BulkOption {
	return func(c *bulkConfig) { c.dryRun = true }
}

// WithDBLister supplies a DBLister used when the source path ends in ".db".
func WithDBLister(lister DBLister) BulkOption {
	return func(c *bulkConfig) { c.dbLister = lister }
}

// withFileLister overrides directory listing (used in tests).
func withFileLister(lister func(string, []string) ([]string, error)) BulkOption {
	return func(c *bulkConfig) { c.fileLister = lister }
}

// ConvertFile converts a single transcript file of any supported format to
// JSON. It returns an *UnknownFileTypeError when the format cannot be
// identified.
func ConvertFile(ctx context.Context, originFile, destinationFile string, metadata map[string]string) error {
	fileType := IdentifyFileType(ctx, originFile)
	if fileType == FileTypeUnknown {
		return newUnknownFileTypeError(originFile)
	}
	return fileConverters()[fileType](originFile, destinationFile, metadata)
}

// job pairs a source file with its converter and destination.
type job struct {
	source      string
	converter   converterFunc
	destination string
}

// BulkConvert converts every transcript under transcriptPath to PodcastIndex
// JSON. transcriptPath may be a directory of transcripts or, when a DBLister is
// configured with WithDBLister, a database (a path ending in ".db"). Conversion
// errors are collected per file rather than aborting the whole run.
func BulkConvert(
	ctx context.Context,
	transcriptPath, destinationPath string,
	opts ...BulkOption,
) (*ConversionSummary, error) {
	cfg := bulkConfig{fileLister: ListFiles}
	for _, opt := range opts {
		opt(&cfg)
	}

	filePaths, metadatas, sourceRoot, err := listSources(ctx, transcriptPath, &cfg)
	if err != nil {
		return nil, err
	}

	grouped := IdentifyFileTypes(ctx, filePaths)
	summary := &ConversionSummary{DryRun: cfg.dryRun, Unknown: grouped[FileTypeUnknown]}

	jobs := buildJobs(grouped, destinationPath, sourceRoot)
	pending := planJobs(jobs, cfg.overwrite, summary)

	if cfg.dryRun {
		for _, j := range pending {
			summary.Converted = append(summary.Converted, ConversionPair{j.source, j.destination})
		}
		return summary, nil
	}
	runJobs(pending, metadatas, summary)
	return summary, nil
}

// listSources resolves the source file paths, per-file metadata, and source
// root for a BulkConvert run.
func listSources(
	ctx context.Context,
	transcriptPath string,
	cfg *bulkConfig,
) (paths []string, metadata map[string]map[string]string, sourceRoot string, err error) {
	if strings.HasSuffix(transcriptPath, ".db") {
		if cfg.dbLister == nil {
			return nil, nil, "", fmt.Errorf("%w for %s", errNoDBLister, transcriptPath)
		}
		paths, metadata, err = cfg.dbLister.ListFiles(ctx, transcriptPath, cfg.ignore)
		return paths, metadata, "", err
	}
	paths, err = cfg.fileLister(transcriptPath, cfg.ignore)
	return paths, map[string]map[string]string{}, transcriptPath, err
}

// buildJobs pairs every known-format file with its converter and unique
// destination, sorted by source path.
func buildJobs(grouped map[FileType][]string, destinationPath, sourceRoot string) []job {
	converters := fileConverters()
	var jobs []job
	for _, fileType := range converterOrder {
		for _, source := range grouped[fileType] {
			jobs = append(jobs, job{source: source, converter: converters[fileType]})
		}
	}
	sort.SliceStable(jobs, func(i, j int) bool { return jobs[i].source < jobs[j].source })

	sources := make([]string, len(jobs))
	for i, j := range jobs {
		sources[i] = j.source
	}
	destinations := assignDestinations(sources, destinationPath, sourceRoot)
	for i := range jobs {
		jobs[i].destination = destinations[jobs[i].source]
	}
	return jobs
}

// planJobs partitions jobs into pending conversions and skipped ones, recording
// the skipped pairs on the summary.
func planJobs(jobs []job, overwrite bool, summary *ConversionSummary) []job {
	var pending []job
	for _, j := range jobs {
		if !overwrite && fileExists(j.destination) {
			summary.Skipped = append(summary.Skipped, ConversionPair{j.source, j.destination})
			continue
		}
		pending = append(pending, j)
	}
	return pending
}

// runJobs converts the pending jobs concurrently and records outcomes on the
// summary, preserving job order.
func runJobs(pending []job, metadatas map[string]map[string]string, summary *ConversionSummary) {
	if len(pending) == 0 {
		return
	}
	results := make([]error, len(pending))
	sem := make(chan struct{}, min(len(pending), maxReadParallelism))
	var wg sync.WaitGroup
	for i := range pending {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			j := pending[i]
			results[i] = runConversion(j.converter, j.source, j.destination, metadatas[j.source])
		}()
	}
	wg.Wait()

	for i, j := range pending {
		if results[i] != nil {
			summary.Failed = append(summary.Failed, ConversionFailure{j.source, results[i].Error()})
		} else {
			summary.Converted = append(summary.Converted, ConversionPair{j.source, j.destination})
		}
	}
}

func runConversion(converter converterFunc, source, destination string, metadata map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(destination), dirPerm); err != nil {
		return err
	}
	return converter(source, destination, metadata)
}

// destinationPath maps a source file to a destination JSON path, mirroring the
// source's directory structure under sourceRoot when possible, otherwise using
// the parent directory name.
func destinationPath(filePath, destinationDir, sourceRoot string) string {
	relative := ""
	if sourceRoot != "" {
		if rel, err := filepath.Rel(sourceRoot, filePath); err == nil &&
			!strings.HasPrefix(rel, "..") && !filepath.IsAbs(rel) {
			relative = rel
		}
	}
	if relative == "" {
		parent := filepath.Base(filepath.Dir(filePath))
		base := filepath.Base(filePath)
		if parent == "." || parent == string(filepath.Separator) || parent == "" {
			relative = base
		} else {
			relative = filepath.Join(parent, base)
		}
	}
	return filepath.Join(destinationDir, withSuffixJSON(relative))
}

// assignDestinations assigns each source a unique destination, deduplicating
// collisions with a " (N)" suffix.
func assignDestinations(filePaths []string, destinationDir, sourceRoot string) map[string]string {
	taken := map[string]struct{}{}
	destinations := make(map[string]string, len(filePaths))
	for _, filePath := range filePaths {
		destination := destinationPath(filePath, destinationDir, sourceRoot)
		candidate := destination
		for counter := 1; ; counter++ {
			if _, ok := taken[candidate]; !ok {
				break
			}
			candidate = dedupeName(destination, counter)
		}
		taken[candidate] = struct{}{}
		destinations[filePath] = candidate
	}
	return destinations
}

// withSuffixJSON replaces a path's final extension with ".json".
func withSuffixJSON(p string) string {
	ext := filepath.Ext(p)
	return p[:len(p)-len(ext)] + ".json"
}

// dedupeName inserts " (counter)" before a destination's extension.
func dedupeName(destination string, counter int) string {
	dir := filepath.Dir(destination)
	base := filepath.Base(destination)
	ext := filepath.Ext(base)
	stem := base[:len(base)-len(ext)]
	return filepath.Join(dir, fmt.Sprintf("%s (%d)%s", stem, counter, ext))
}

// fileExists reports whether path exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
