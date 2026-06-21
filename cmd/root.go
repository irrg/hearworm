package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hearworm",
	Short: "Audiobook conversion and manipulation tool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		for _, bin := range []string{"ffmpeg", "ffprobe"} {
			if _, err := exec.LookPath(bin); err != nil {
				return fmt.Errorf("%s not found in PATH — install ffmpeg >= 4.1", bin)
			}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
