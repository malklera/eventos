// Package cmd contains all logic regarding the eventos command
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eventos",
	Short: "Suit to manage videos from 360 platform.",
	Long: `Steps to follow.

	eventos init eventName

Create Gmail Drive folder for the event, ensure it is private, copy <url>	

	cd eventName
	eventos qr <url>

Copy videos from phone to eventName/original

	eventos copy <path to phone videos directory> <event-name>

	eventos rename

	eventos format

Cut each video manually and export them to ~/Videos/eventos/event-name/cortado
In the process take a screenshot of each for the catalog.

Put all images into the "catalogo" directory.

	mv ~/Videos/eventos/<event-name>/cortado/*.png ~/Videos/eventos/<event-name>/catalogo
	cd ~/Videos/eventos/<event-name>/catalogo
	perl-rename -ni 's/^(\d+)-.*/$1.png/' *.png

See the output of that, if it look ok, remove the -n flag

eventos edit...
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.eventos.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
