package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"hearworm/internal/chapter"
)

type ProbeResult struct {
	Format   FormatInfo   `json:"format"`
	Streams  []StreamInfo `json:"streams"`
	Chapters []FFChapter  `json:"chapters"`
}

type FormatInfo struct {
	Filename string            `json:"filename"`
	Duration string            `json:"duration"`
	Tags     map[string]string `json:"tags"`
}

type StreamInfo struct {
	CodecType  string `json:"codec_type"`
	CodecName  string `json:"codec_name"`
	SampleRate string `json:"sample_rate"`
	Channels   int    `json:"channels"`
	BitRate    string `json:"bit_rate"`
}

type FFChapter struct {
	ID        int               `json:"id"`
	TimeBase  string            `json:"time_base"`
	Start     int64             `json:"start"`
	StartTime string            `json:"start_time"`
	End       int64             `json:"end"`
	EndTime   string            `json:"end_time"`
	Tags      map[string]string `json:"tags"`
}

func Probe(path string) (*ProbeResult, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_chapters",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed for %q: %w", path, err)
	}
	var result ProbeResult
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("ffprobe output parse failed: %w", err)
	}
	return &result, nil
}

// ParseSeconds parses a float seconds string (ffprobe format) into time.Duration.
func ParseSeconds(s string) (time.Duration, error) {
	var secs float64
	if _, err := fmt.Sscanf(s, "%f", &secs); err != nil {
		return 0, fmt.Errorf("cannot parse %q as seconds: %w", s, err)
	}
	return time.Duration(secs * float64(time.Second)), nil
}

// ChaptersFromProbe converts ffprobe chapter data to a chapter.List.
// Uses StartTime/EndTime (float second strings) for accuracy across all timebases.
func ChaptersFromProbe(r *ProbeResult) chapter.List {
	list := make(chapter.List, len(r.Chapters))
	for i, c := range r.Chapters {
		var startSecs, endSecs float64
		fmt.Sscanf(c.StartTime, "%f", &startSecs)
		fmt.Sscanf(c.EndTime, "%f", &endSecs)
		list[i] = chapter.Chapter{
			Title: c.Tags["title"],
			Start: time.Duration(startSecs * float64(time.Second)),
			End:   time.Duration(endSecs * float64(time.Second)),
		}
	}
	return list
}
