package cmd

import (
	"fmt"
	"time"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy files from phone to pc.",
	Long:  `eventos copy srcPath dstPath
	eventos copy ".../Almacenamiento interno/DCIM/Camera/<eventName>" <eventName>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("wrong number of arguments, got: %d, want: 2", len(args))
		}

		files, err := listVideos(args[0])
		if err != nil {
			return fmt.Errorf("error reading '%s': %w", args[0], err)
		}

		fmt.Println("Total src videos:", len(files))
		count := 0
		evento := filepath.Join("/home/malklera/Videos/eventos", args[1])
		original := filepath.Join(evento, "original")
		start := time.Now()
		for _, file := range files {
			src := filepath.Join(args[0], file.Name())
			dst := filepath.Join(original, file.Name())
			err := copyFile(src, dst)
			if err != nil {
				fmt.Fprintf(os.Stderr, "copyFile(%s, %s): %v\n", src, dst, err)
			} else {
				count++
				fmt.Println("Copied:", src)
			}
		}
		fmt.Println("Total copied:", count)
		fmt.Println("Time taken:", time.Since(start))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
