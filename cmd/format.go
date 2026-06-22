package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Makes the videos fixed frame rate and took out sound.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := "renombrado"
		dst := "formateado"
		files, err := os.ReadDir(src)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", src, err)
		}
		count := len(files)
		done := 0
		start := time.Now()
		for _, file := range files {
			srcF := filepath.Join(src, file.Name())
			dstF := filepath.Join(dst, file.Name())
			fmt.Println("Formateando:", srcF)
			ffmpeg := exec.Command("ffmpeg",
				"-v",
				"warning",
				"-i",
				srcF,
				"-vf",
				"fps=30",
				"-c:v",
				"libx264",
				"-preset",
				"medium",
				"-crf",
				"23",
				"-movflags",
				"+faststart",
				"-an",
				dstF)
			ffmpeg.Stdout = os.Stdout
			ffmpeg.Stderr = os.Stderr
			err := ffmpeg.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error running %s: %v\n", ffmpeg.String(), err)
			} else {
				done++
			}
		}
		fmt.Println("Total files:", count)
		fmt.Println("Files formated:", done)
		fmt.Println("Time taken:", time.Since(start))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
}
