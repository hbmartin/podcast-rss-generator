package overcastdb_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/hbmartin/podcast-rss-generator/v2/transcript/overcastdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE episodes_extended (
    transcriptDownloadPath TEXT,
    title TEXT,
    enclosureUrl TEXT,
    guid TEXT,
    feedXmlUrl TEXT
);
CREATE TABLE feeds_extended (
    title TEXT,
    xmlUrl TEXT
);
INSERT INTO feeds_extended VALUES
    ('My Show', 'https://example.com/feed.xml');
INSERT INTO episodes_extended VALUES
    ('/transcripts/ep1.srt', 'Episode 1', 'https://example.com/ep1.mp3',
     'guid-1', 'https://example.com/feed.xml'),
    ('/transcripts/ignored.srt', 'Episode 2', 'https://example.com/ep2.mp3',
     'guid-2', 'https://example.com/feed.xml'),
    (NULL, 'No transcript', 'https://example.com/ep3.mp3',
     'guid-3', 'https://example.com/feed.xml');
`

func createDB(t *testing.T) string {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "overcast.db")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()
	_, err = db.ExecContext(context.Background(), schema)
	require.NoError(t, err)
	return dbPath
}

func TestListFiles(t *testing.T) {
	t.Parallel()
	dbPath := createDB(t)

	files, metadata, err := overcastdb.ListFiles(context.Background(), dbPath, []string{"ignored.srt"})
	require.NoError(t, err)
	assert.Equal(t, []string{"/transcripts/ep1.srt"}, files)

	got := metadata["/transcripts/ep1.srt"]
	assert.Equal(t, "Episode 1", got["title"])
	assert.Equal(t, "https://example.com/ep1.mp3", got["enclosureUrl"])
	assert.Equal(t, "guid-1", got["guid"])
	assert.Equal(t, "My Show", got["feedTitle"])
	assert.Equal(t, "https://example.com/feed.xml", got["xmlUrl"])
}

func TestListFilesImplementsDBLister(t *testing.T) {
	t.Parallel()
	dbPath := createDB(t)

	files, metadata, err := overcastdb.Lister{}.ListFiles(context.Background(), dbPath, nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"/transcripts/ep1.srt", "/transcripts/ignored.srt"}, files)
	assert.Len(t, metadata, 2)
}
