package ffmpeg_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"hearworm/internal/chapter"
	"hearworm/internal/ffmpeg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteMeta_Content(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "meta.txt")

	tags := map[string]string{"title": "Test Book", "artist": "Test Author"}
	chapters := chapter.List{
		{Title: "Intro", Start: 0, End: 5 * time.Second},
		{Title: "Chapter 1", Start: 5 * time.Second, End: 65 * time.Second},
	}

	err := ffmpeg.WriteMeta(path, tags, chapters, 65*time.Second)
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	content := string(data)

	assert.Contains(t, content, ";FFMETADATA1")
	assert.Contains(t, content, "title=Test Book")
	assert.Contains(t, content, "[CHAPTER]")
	assert.Contains(t, content, "TIMEBASE=1/1000")
	assert.Contains(t, content, "title=Intro")
	assert.Contains(t, content, "START=0")
	assert.Contains(t, content, "END=5000")
	assert.Contains(t, content, "title=Chapter 1")
	assert.Contains(t, content, "START=5000")
	assert.Contains(t, content, "END=65000")
}
