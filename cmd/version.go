package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
)

var tag = ""

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version and build information.",
	RunE: func(cmd *cobra.Command, args []string) error {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			return fmt.Errorf("failed to debug.ReadBuildInfo()")
		}
		fmt.Print("eventos ")
		fmt.Printf("version %s ", tag)
		vcsTime := ""
		for _, s := range bi.Settings {
			if s.Key == "vcs.time" {
				vcsTime = s.Value
			}
		}
		date, err := time.Parse(time.RFC3339, vcsTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nerror parsing time: %v\n", err)
		}
		fmt.Printf("(%v) ", date.Format(time.DateOnly))
		fmt.Println("copyright 2025 Malklera")
		fmt.Println()
		fmt.Println("Source", bi.Path)
		fmt.Println()
		fmt.Println("Build information.")
		fmt.Println(bi.GoVersion)
		for _, s := range bi.Settings {
			fmt.Printf("%s=%s\n", s.Key, s.Value)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
