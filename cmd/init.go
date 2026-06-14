/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	// TODO: add the target flag to operate in other place instead of where the
	// command is run
	Use:   "init",
	Short: "Create the required directories to contain the event",
	Long:  `e.g. ./eventos init eventName`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		event := args[0]
		err := os.Mkdir(event, 0750)
		if err != nil {
			fmt.Println("error:", err)
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
				return err
			}
		}
		return nil
	},
}

// init creates the directories needed for the event
func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
