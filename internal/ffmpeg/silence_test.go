package ffmpeg_test

import (
	"testing"
	"time"

	"hearworm/internal/ffmpeg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSilenceOutput(t *testing.T) {
	stderr := `[silencedetect @ 0x7f] silence_start: 5.5
[silencedetect @ 0x7f] silence_end: 6.2 | silence_duration: 0.7
[silencedetect @ 0x7f] silence_start: 30
[silencedetect @ 0x7f] silence_end: 31.5 | silence_duration: 1.5`

	ranges, err := ffmpeg.ParseSilenceOutput(stderr)
	require.NoError(t, err)
	require.Len(t, ranges, 2)
	assert.Equal(t, 5500*time.Millisecond, ranges[0].Start)
	assert.Equal(t, 6200*time.Millisecond, ranges[0].End)
	assert.Equal(t, 30000*time.Millisecond, ranges[1].Start)
	assert.Equal(t, 31500*time.Millisecond, ranges[1].End)
}

func TestParseSilenceOutput_Empty(t *testing.T) {
	ranges, err := ffmpeg.ParseSilenceOutput("nothing here")
	require.NoError(t, err)
	assert.Empty(t, ranges)
}
