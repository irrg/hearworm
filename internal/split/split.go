package split

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"m4b/internal/ffmpeg"
)

type Options struct {
	InputFile string
	OutputDir string
	Format    string // "m4b", "mp3", "m4a"
}

func Run(opts Options) error {
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return err
	}
	probe, err := ffmpeg.Probe(opts.InputFile)
	if err != nil {
		return fmt.Errorf("probe %q: %w", opts.InputFile, err)
	}
	chapters := ffmpeg.ChaptersFromProbe(probe)
	if len(chapters) == 0 {
		return fmt.Errorf("no chapters found in %q", opts.InputFile)
	}

	totalDur, _ := ffmpeg.ParseSeconds(probe.Format.Duration)
	ext := "." + coalesce(opts.Format, "m4b")

	for i, ch := range chapters {
		end := ch.End
		if end == 0 {
			if i+1 < len(chapters) {
				end = chapters[i+1].Start
			} else {
				end = totalDur
			}
		}

		outFile := filepath.Join(opts.OutputDir,
			fmt.Sprintf("%03d-%s%s", i+1, sanitize(ch.Title), ext))

		metaFlags := []string{
			fmt.Sprintf("title=%s", ch.Title),
			fmt.Sprintf("track=%d/%d", i+1, len(chapters)),
		}
		for k, v := range probe.Format.Tags {
			if k != "title" && v != "" {
				metaFlags = append(metaFlags, fmt.Sprintf("%s=%s", k, v))
			}
		}

		if err := ffmpeg.ExtractSegment(opts.InputFile, outFile, ch.Start, end, metaFlags); err != nil {
			return fmt.Errorf("extract chapter %d %q: %w", i+1, ch.Title, err)
		}
		fmt.Printf("[%d/%d] %s\n", i+1, len(chapters), ch.Title)
	}
	return nil
}

func sanitize(s string) string {
	r := strings.NewReplacer(
		"/", "-", "\\", "-", ":", "-", "*", "",
		"?", "", "\"", "", "<", "", ">", "", "|", "",
	)
	return r.Replace(s)
}

func coalesce(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
