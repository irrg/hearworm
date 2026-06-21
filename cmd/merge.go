package cmd

import (
	"github.com/spf13/cobra"
	"hearworm/internal/merge"
)

var mergeCmd = &cobra.Command{
	Use:   "merge <input-dir>",
	Short: "Merge audio files in a directory into a single m4b",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		var opts merge.Options
		opts.InputDir = args[0]
		opts.OutputFile, _ = f.GetString("output-file")
		opts.Tag.Title, _ = f.GetString("name")
		opts.Tag.Artist, _ = f.GetString("artist")
		opts.Tag.Album, _ = f.GetString("album")
		opts.Tag.AlbumArtist, _ = f.GetString("album-artist")
		opts.Tag.Genre, _ = f.GetString("genre")
		opts.Tag.Year, _ = f.GetString("year")
		opts.Tag.Series, _ = f.GetString("series")
		opts.Tag.SeriesPart, _ = f.GetString("series-part")
		opts.Tag.Bitrate, _ = f.GetString("audio-bitrate")
		opts.Format, _ = f.GetString("audio-format")
		return merge.Run(opts)
	},
}

func init() {
	f := mergeCmd.Flags()
	f.StringP("output-file", "o", "", "output .m4b path (required)")
	f.String("name", "", "audiobook title")
	f.String("artist", "", "author/artist")
	f.String("album", "", "album (defaults to --name)")
	f.String("album-artist", "", "album artist tag")
	f.String("genre", "Audiobook", "genre tag")
	f.String("year", "", "publication year")
	f.String("series", "", "series name (used to compute sort order)")
	f.String("series-part", "", "series part number")
	f.String("audio-format", "m4b", "output format: m4b, m4a, mp3")
	f.String("audio-bitrate", "", "output audio bitrate (e.g. 64k, 128k); defaults to 128k for mp3, 64k for aac")
	mergeCmd.MarkFlagRequired("output-file")
	rootCmd.AddCommand(mergeCmd)
}
