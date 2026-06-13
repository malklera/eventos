/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
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
	Run: func(cmd *cobra.Command, args []string) {
		files, err := os.ReadDir("original")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading 'original': %v\n", err)
			os.Exit(1)
		}

		for n, file := range files {
			src := filepath.Join("original", file.Name())
			dst := filepath.Join("renombrado", strconv.Itoa(n)+".mp4")
			srcF, err := os.Open(src)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening '%s': %v\n", src, err)
			}
			defer srcF.Close()

			dstF, err := os.Create(dst)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating '%s': %v\n", src, err)
			}
			defer dstF.Close()

			_, err = io.Copy(dstF, srcF)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error copying '%s' to '%s': %v\n", src, dst, err)
			}
			if err := dstF.Sync(); err != nil {
				fmt.Fprintf(os.Stderr, "error dstF.Sync(): %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
