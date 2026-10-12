package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// catalogCmd represents the catalog command
var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Copy the screenshots from <nameEvent>/cortado to <nameEvent>/catalogo",
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := os.ReadDir("cortado")
		if err != nil {
			return fmt.Errorf("os.ReadDir(\"cortado\"): %w", err)
		}
		img := 0
		vid := 0
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".png") {
				img += 1
				number, _, _ := strings.Cut(file.Name(), "-")
				src := filepath.Join("cortado", file.Name())
				dst := filepath.Join("catalogo", number+".png")
				err := os.Rename(src, dst)
				if err != nil {
					fmt.Fprintf(os.Stderr, "os.Rename(%s, %s): %v\n", src, dst, err)
				}
			} else {
				vid += 1
			}
		}

		fmt.Println("Images:", img)
		fmt.Println("Videos:", vid)

		fp, err := os.ReadDir("catalogo")
		if err != nil {
			return fmt.Errorf("os.ReadDir(\"catalogo\"): %w", err)
		}
		fmt.Println("Catalog:", len(fp))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catalogCmd)
}
