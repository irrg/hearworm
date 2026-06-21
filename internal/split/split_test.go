package split_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"hearworm/internal/split"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeChapteredM4b(t *testing.T, path string) {
	t.Helper()
	dir := filepath.Dir(path)
	audio := filepath.Join(dir, "raw.m4a")
	meta := filepath.Join(dir, "meta.txt")

	cmd := exec.Command("ffmpeg", "-y", "-f", "lavfi",
		"-i", "sine=frequency=440:duration=10",
		"-c:a", "aac", "-b:a", "64k", audio)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	metaContent := ";FFMETADATA1\ntitle=Test\n\n" +
		"[CHAPTER]\nTIMEBASE=1/1000\nSTART=0\nEND=5000\ntitle=Part One\n\n" +
		"[CHAPTER]\nTIMEBASE=1/1000\nSTART=5000\nEND=10000\ntitle=Part Two\n\n"
	require.NoError(t, os.WriteFile(meta, []byte(metaContent), 0644))

	cmd = exec.Command("ffmpeg", "-y", "-i", audio, "-i", meta,
		"-map_metadata", "1", "-c", "copy", path)
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))
}

func TestSplit_TwoChapters(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available")
	}
	dir := t.TempDir()
	input := filepath.Join(dir, "input.m4b")
	outDir := filepath.Join(dir, "output")
	makeChapteredM4b(t, input)

	err := split.Run(split.Options{
		InputFile: input,
		OutputDir: outDir,
		Format:    "m4b",
	})
	require.NoError(t, err)

	entries, err := os.ReadDir(outDir)
	require.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "001-Part One.m4b", entries[0].Name())
	assert.Equal(t, "002-Part Two.m4b", entries[1].Name())
}

func TestSplit_NoChapters(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available")
	}
	dir := t.TempDir()
	input := filepath.Join(dir, "input.m4a")
	cmd := exec.Command("ffmpeg", "-y", "-f", "lavfi",
		"-i", "sine=frequency=440:duration=5",
		"-c:a", "aac", "-b:a", "64k", input)
	require.NoError(t, cmd.Run())

	err := split.Run(split.Options{InputFile: input, OutputDir: t.TempDir(), Format: "m4b"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no chapters")
}
