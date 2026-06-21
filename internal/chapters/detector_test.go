package chapters_test

import (
	"testing"
	"time"

	"hearworm/internal/chapters"
	"hearworm/internal/ffmpeg"
	"github.com/stretchr/testify/assert"
)

func TestSilencesToChapters_TwoSilences(t *testing.T) {
	silences := []ffmpeg.SilenceRange{
		{Start: 5 * time.Second, End: 6 * time.Second},
		{Start: 30 * time.Second, End: 31 * time.Second},
	}
	list := chapters.SilencesToChapters(silences, 60*time.Second)
	assert.Len(t, list, 3)
	assert.Equal(t, time.Duration(0), list[0].Start)
	assert.Equal(t, 5500*time.Millisecond, list[0].End)
	assert.Equal(t, 5500*time.Millisecond, list[1].Start)
	assert.Equal(t, 30500*time.Millisecond, list[1].End)
	assert.Equal(t, 30500*time.Millisecond, list[2].Start)
	assert.Equal(t, 60*time.Second, list[2].End)
	assert.Equal(t, "Chapter 1", list[0].Title)
	assert.Equal(t, "Chapter 3", list[2].Title)
}

func TestSilencesToChapters_NoSilences(t *testing.T) {
	list := chapters.SilencesToChapters(nil, 120*time.Second)
	assert.Len(t, list, 1)
	assert.Equal(t, "Chapter 1", list[0].Title)
	assert.Equal(t, 120*time.Second, list[0].End)
}
