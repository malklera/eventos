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
		fmt.Println("catalog called")
		files, err := os.ReadDir("cortado")
		if err != nil {
			return fmt.Errorf("os.ReadDir(\"cortado\"): %w", err)
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".png") {
				number, _, _ := strings.Cut(file.Name(), "-")
				err := os.Rename(filepath.Join("cortado", file.Name()), filepath.Join("catalogo", number+".png"))
				if err != nil {
					fmt.Fprintf(os.Stderr, "os.Rename(%s, %s+\".png\"): %v\n", file.Name(), number, err)
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catalogCmd)
}
