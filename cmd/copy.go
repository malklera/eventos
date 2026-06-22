package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy files from phone to pc.",
	Long:  `eventos copy srcPath dstPath
	eventos copy .../Almacenamiento interno/DCIM/Camera/<eventName> <eventName>`,
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
		originals := filepath.Join("~/Videos/eventos", args[1])
		for _, file := range files {
			src := filepath.Join(args[0], file.Name())
			dst := filepath.Join(originals, file.Name())
			err := copyFile(src, dst)
			if err != nil {
				fmt.Fprintf(os.Stderr, "copyFile(%s, %s): %v\n", src, dst, err)
			} else {
				count++
			}
		}
		fmt.Println("Total copied:", count)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
