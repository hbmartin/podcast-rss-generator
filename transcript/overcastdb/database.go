// Package overcastdb reads transcript file paths and episode metadata from an
// overcast-to-sqlite database. It is a separate package so that consumers of the
// transcript converters do not pull in the pure-Go SQLite driver unless they
// use this database integration.
//
// https://github.com/hbmartin/overcast-to-sqlite
package overcastdb

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	_ "modernc.org/sqlite" // pure-Go SQLite driver registered as "sqlite".
)

// Metadata keys attached to each transcript, matching the overcast-to-sqlite
// column names.
const (
	enclosureURL = "enclosureUrl"
	feedTitle    = "feedTitle"
	guid         = "guid"
	title        = "title"
	xmlURL       = "xmlUrl"
)

// selectQuery selects transcript download paths and episode/feed metadata for
// every episode that has a transcript.
const selectQuery = `SELECT transcriptDownloadPath, episodes_extended.title, ` +
	`episodes_extended.enclosureUrl, episodes_extended.guid, ` +
	`feeds_extended.title, feeds_extended.xmlUrl ` +
	`FROM episodes_extended ` +
	`LEFT JOIN feeds_extended ON episodes_extended.feedXmlUrl = feeds_extended.xmlUrl ` +
	`WHERE transcriptDownloadPath IS NOT NULL`

// Lister reads transcripts from an overcast-to-sqlite database. It implements
// transcript.DBLister so it can be passed to transcript.BulkConvert via
// transcript.WithDBLister.
type Lister struct{}

// ListFiles implements transcript.DBLister.
func (Lister) ListFiles(
	ctx context.Context,
	dbPath string,
	ignore []string,
) ([]string, map[string]map[string]string, error) {
	return ListFiles(ctx, dbPath, ignore)
}

// ListFiles returns the transcript file paths (excluding ignored and disallowed
// names) and their episode metadata from the overcast-to-sqlite database at
// dbPath.
func ListFiles(
	ctx context.Context,
	dbPath string,
	ignore []string,
) (files []string, metadata map[string]map[string]string, err error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("opening database %s: %w", dbPath, err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("closing database: %w", closeErr)
		}
	}()

	rows, err := db.QueryContext(ctx, selectQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("querying transcripts: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("closing rows: %w", closeErr)
		}
	}()

	metadata = map[string]map[string]string{}
	for rows.Next() {
		var path, episodeTitle, url, episodeGUID, feed, feedXML sql.NullString
		if err = rows.Scan(&path, &episodeTitle, &url, &episodeGUID, &feed, &feedXML); err != nil {
			return nil, nil, fmt.Errorf("scanning row: %w", err)
		}
		if !path.Valid || !isFileAllowed(baseName(path.String), ignore) {
			continue
		}
		files = append(files, path.String)
		metadata[path.String] = map[string]string{
			enclosureURL: url.String,
			guid:         episodeGUID.String,
			title:        episodeTitle.String,
			feedTitle:    feed.String,
			xmlURL:       feedXML.String,
		}
	}
	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterating rows: %w", err)
	}
	return files, metadata, nil
}

// isFileAllowed mirrors the transcript package's filter: the name must not be
// ignored, hidden, or a PDF or octet-stream.
func isFileAllowed(filename string, ignore []string) bool {
	if slices.Contains(ignore, filename) {
		return false
	}
	return !strings.HasPrefix(filename, ".") &&
		!strings.HasSuffix(filename, ".pdf") &&
		!strings.HasSuffix(filename, ".octet-stream")
}

// baseName returns the final "/"-separated component of a stored path.
func baseName(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}
