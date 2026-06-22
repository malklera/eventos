package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename the files into the order they where created.",
	Long:  `I say rename, but it is a copy with a different final name.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := listVideos("original")
		if err != nil {
			return fmt.Errorf("error reading 'original': %w", err)
		}

		for n, file := range files {
			src := filepath.Join("original", file.Name())
			dst := filepath.Join("renombrado", strconv.Itoa(n+1)+".mp4")
			err := copyFile(src, dst)
			if err != nil {
				fmt.Fprintf(os.Stderr, "copyFile(%s, %s): %v\n", src, dst, err)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
