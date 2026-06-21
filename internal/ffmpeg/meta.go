package ffmpeg

import (
	"fmt"
	"os"
	"strings"
	"time"

	"hearworm/internal/chapter"
)

func WriteMeta(path string, tags map[string]string, chapters chapter.List, totalDuration time.Duration) error {
	var sb strings.Builder
	sb.WriteString(";FFMETADATA1\n")
	for k, v := range tags {
		if v != "" {
			fmt.Fprintf(&sb, "%s=%s\n", k, escapeMeta(v))
		}
	}
	sb.WriteString("\n")
	for _, ch := range chapters {
		endMs := ch.End.Milliseconds()
		if endMs == 0 {
			endMs = totalDuration.Milliseconds()
		}
		fmt.Fprintf(&sb, "[CHAPTER]\nTIMEBASE=1/1000\nSTART=%d\nEND=%d\ntitle=%s\n\n",
			ch.Start.Milliseconds(), endMs, escapeMeta(ch.Title))
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func escapeMeta(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "=", "\\=")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, "#", "\\#")
	s = strings.ReplaceAll(s, "\n", "\\\n")
	return s
}
