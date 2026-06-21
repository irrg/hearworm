package cmd

import (
	"github.com/spf13/cobra"
	"m4b/internal/split"
)

var splitCmd = &cobra.Command{
	Use:   "split <input.m4b>",
	Short: "Split an m4b into one file per chapter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		outDir, _ := f.GetString("output-dir")
		format, _ := f.GetString("audio-format")
		return split.Run(split.Options{
			InputFile: args[0],
			OutputDir: outDir,
			Format:    format,
		})
	},
}

func init() {
	f := splitCmd.Flags()
	f.StringP("output-dir", "o", "", "directory to write chapter files (required)")
	f.String("audio-format", "m4b", "output format: m4b, mp3, m4a")
	splitCmd.MarkFlagRequired("output-dir")
	rootCmd.AddCommand(splitCmd)
}
