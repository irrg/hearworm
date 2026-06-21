# m4b-go

Convert, split, merge, and chapter-mark audiobook files. Single static binary. **ffmpeg is the only runtime dependency.**

Supports mp3, m4b, m4a, aac, flac, ogg, wav, and wma as input or output. The primary use case is splitting m4b audiobooks into mp3 files for devices that don't support the m4b format.

## License

MIT License — see [LICENSE](LICENSE) for full text.

## Requirements

- Go 1.22+ (to build)
- ffmpeg ≥ 4.1 (runtime, must be in PATH)

## Install

```bash
git clone https://github.com/irrg/m4b-go
cd m4b-go
go build -o m4b-tool .
```

## Usage

### split

Split an m4b into one mp3 (or other format) per chapter.

```bash
m4b-tool split input.m4b --output-dir ./chapters/ --audio-format mp3
```

Output files are named `001-Chapter Title.mp3`. Tags from the source file are preserved on each output file. If source audio is already the target codec, it is remuxed without re-encoding.

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output-dir` | *(required)* | Directory for output files |
| `--audio-format` | `m4b` | Output format: `mp3`, `m4b`, `m4a` |

### merge

Combine a directory of audio files into a single file with embedded chapters.

```bash
m4b-tool merge /path/to/audio-files/ \
  --output-file output.m4b \
  --name "The Name of the Wind" \
  --artist "Patrick Rothfuss" \
  --series "Kingkiller Chronicle" \
  --series-part "1"
```

One chapter is created per input file, named from the file's `title` tag (or the filename if no tag exists). Files are sorted alphabetically — prefix with track numbers (`01-`, `02-`, etc.) for correct order. If source audio already matches the output codec, it is remuxed without re-encoding.

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output-file` | *(required)* | Output file path |
| `--audio-format` | `m4b` | Output format: `mp3`, `m4b`, `m4a` |
| `--audio-bitrate` | `128k` mp3 / `64k` aac | Output bitrate |
| `--name` | | Audiobook title |
| `--artist` | | Author / artist |
| `--album` | same as `--name` | Album tag |
| `--album-artist` | | Album artist tag |
| `--genre` | `Audiobook` | Genre tag |
| `--year` | | Publication year |
| `--series` | | Series name (used to compute sort order) |
| `--series-part` | | Series part number |

**Series sort order:** When `--series` and `--series-part` are set, `sort_name` is automatically computed as `"{series} {part} - {title}"` (e.g. `Kingkiller Chronicle 1 - The Name of the Wind`). Players that respect sort tags will order the series correctly regardless of alphabetical title order.

### chapters

Detect silence in an audio file and embed chapter markers at silence midpoints.

```bash
m4b-tool chapters input.m4b --output-file output.m4b
```

Omit `--output-file` to overwrite the input in-place.

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output-file` | *(overwrites input)* | Output file path |
| `--silence-min-noise` | `-30` | Noise threshold in dB |
| `--silence-min-length` | `0.5` | Minimum silence duration in seconds |
