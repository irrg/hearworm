package ffmpeg

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SilenceRange struct {
	Start time.Duration
	End   time.Duration
}

var silenceStartRe = regexp.MustCompile(`silence_start: ([0-9.]+)`)
var silenceEndRe = regexp.MustCompile(`silence_end: ([0-9.]+)`)

func DetectSilence(path string, noiseDb float64, minDuration float64) ([]SilenceRange, error) {
	filter := fmt.Sprintf("silencedetect=noise=%.0fdB:duration=%.2f", noiseDb, minDuration)
	cmd := exec.Command("ffmpeg", "-i", path, "-af", filter, "-f", "null", "-")
	out, _ := cmd.CombinedOutput() // ffmpeg exits non-zero for -f null; ignore exit code
	return ParseSilenceOutput(string(out))
}

func ParseSilenceOutput(output string) ([]SilenceRange, error) {
	var ranges []SilenceRange
	var pending *SilenceRange
	for _, line := range strings.Split(output, "\n") {
		if m := silenceStartRe.FindStringSubmatch(line); m != nil {
			secs, err := strconv.ParseFloat(m[1], 64)
			if err != nil {
				return nil, fmt.Errorf("parse silence_start %q: %w", m[1], err)
			}
			pending = &SilenceRange{Start: floatToMs(secs)}
		}
		if m := silenceEndRe.FindStringSubmatch(line); m != nil && pending != nil {
			secs, err := strconv.ParseFloat(m[1], 64)
			if err != nil {
				return nil, fmt.Errorf("parse silence_end %q: %w", m[1], err)
			}
			pending.End = floatToMs(secs)
			ranges = append(ranges, *pending)
			pending = nil
		}
	}
	return ranges, nil
}

func floatToMs(secs float64) time.Duration {
	return time.Duration(int64(secs*1000)) * time.Millisecond
}
