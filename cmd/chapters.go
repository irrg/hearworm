package cmd

import (
	"github.com/spf13/cobra"
	"hearworm/internal/chapters"
)

var chaptersCmd = &cobra.Command{
	Use:   "chapters <input.m4b>",
	Short: "Detect silence and add chapter markers to an m4b",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		noise, _ := f.GetFloat64("silence-min-noise")
		dur, _ := f.GetFloat64("silence-min-length")
		out, _ := f.GetString("output-file")
		return chapters.Run(chapters.Options{
			InputFile:  args[0],
			OutputFile: out,
			SilenceDb:  noise,
			SilenceDur: dur,
		})
	},
}

func init() {
	f := chaptersCmd.Flags()
	f.StringP("output-file", "o", "", "output file (default: overwrite input)")
	f.Float64("silence-min-noise", -30, "silence noise threshold in dB (e.g. -30)")
	f.Float64("silence-min-length", 0.5, "minimum silence duration in seconds to split on")
	rootCmd.AddCommand(chaptersCmd)
}
