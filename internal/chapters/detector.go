package chapters

import (
	"fmt"
	"os"
	"time"

	"m4b/internal/chapter"
	"m4b/internal/ffmpeg"
)

type Options struct {
	InputFile  string
	OutputFile string  // if empty, overwrites InputFile
	SilenceDb  float64 // default -30
	SilenceDur float64 // default 0.5 seconds
}

func SilencesToChapters(silences []ffmpeg.SilenceRange, totalDur time.Duration) chapter.List {
	if len(silences) == 0 {
		return chapter.List{{Title: "Chapter 1", Start: 0, End: totalDur}}
	}
	var boundaries []time.Duration
	for _, s := range silences {
		boundaries = append(boundaries, s.Start+(s.End-s.Start)/2)
	}
	var list chapter.List
	prev := time.Duration(0)
	for i, b := range boundaries {
		list = append(list, chapter.Chapter{
			Title: fmt.Sprintf("Chapter %d", i+1),
			Start: prev,
			End:   b,
		})
		prev = b
	}
	list = append(list, chapter.Chapter{
		Title: fmt.Sprintf("Chapter %d", len(boundaries)+1),
		Start: prev,
		End:   totalDur,
	})
	return list
}

func Run(opts Options) error {
	if opts.SilenceDb == 0 {
		opts.SilenceDb = -30
	}
	if opts.SilenceDur == 0 {
		opts.SilenceDur = 0.5
	}
	outputFile := opts.OutputFile
	if outputFile == "" {
		outputFile = opts.InputFile
	}

	probe, err := ffmpeg.Probe(opts.InputFile)
	if err != nil {
		return fmt.Errorf("probe: %w", err)
	}
	totalDur, err := ffmpeg.ParseSeconds(probe.Format.Duration)
	if err != nil {
		return err
	}

	silences, err := ffmpeg.DetectSilence(opts.InputFile, opts.SilenceDb, opts.SilenceDur)
	if err != nil {
		return fmt.Errorf("silence detection: %w", err)
	}
	fmt.Printf("found %d silence(s)\n", len(silences))

	chs := SilencesToChapters(silences, totalDur)
	fmt.Printf("generated %d chapter(s)\n", len(chs))

	tags := map[string]string{}
	for k, v := range probe.Format.Tags {
		tags[k] = v
	}

	metaFile, err := os.CreateTemp("", "m4b-meta-*.txt")
	if err != nil {
		return err
	}
	metaFile.Close()
	defer os.Remove(metaFile.Name())

	if err := ffmpeg.WriteMeta(metaFile.Name(), tags, chs, totalDur); err != nil {
		return err
	}

	// Write to temp first when overwriting in-place to avoid clobbering source
	if outputFile == opts.InputFile {
		tmp, err := os.CreateTemp("", "m4b-out-*.m4b")
		if err != nil {
			return err
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		if err := ffmpeg.ApplyMeta(opts.InputFile, tmp.Name(), metaFile.Name()); err != nil {
			return err
		}
		return os.Rename(tmp.Name(), outputFile)
	}

	return ffmpeg.ApplyMeta(opts.InputFile, outputFile, metaFile.Name())
}
