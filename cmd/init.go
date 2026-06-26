package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create the required directories to contain the event",
	Long: `e.g. eventos init eventName
	eventos init ~/Videos/eventos/eventName`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		event := args[0]
		err := os.Mkdir(event, 0750)
		if err != nil {
			return fmt.Errorf("failed to init: %w", err)
		}
		err2 := os.Mkdir(filepath.Join(event, "original"), 0750)
		err3 := os.Mkdir(filepath.Join(event, "renombrado"), 0750)
		err4 := os.Mkdir(filepath.Join(event, "formateado"), 0750)
		err5 := os.Mkdir(filepath.Join(event, "cortado"), 0750)
		err6 := os.Mkdir(filepath.Join(event, "editado"), 0750)
		err7 := os.Mkdir(filepath.Join(event, "catalogo"), 0750)
		if err2 != nil ||
			err3 != nil ||
			err4 != nil ||
			err5 != nil ||
			err6 != nil ||
			err7 != nil {
			if err := os.RemoveAll(event); err != nil {
				return fmt.Errorf("failed to init: %w", err)
			}
		}
		return nil
	},
}

// init creates the directories needed for the event
func init() {
	rootCmd.AddCommand(initCmd)
}
