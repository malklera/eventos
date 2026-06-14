/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

var (
	logo        string
	music       string
	image       string
	eventText   string
	clientText  string
	up          string
	down        string
	clientColor string
	textColor   string
	font        string
	leftIcon    string
	rightIcon   string
	save        bool
	run         bool
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Command to edit the videos",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		videoWidth := 1080
		videoHeight := 1920
		transitionDuration := 1
		clientImgTime := 2

		// Text styling
		lineSpacing := 10
		fontSize := 80
		borderWidth := 2
		borderColor := "#000000"
		boxPadding := 8
		clientTextX := "(w-tw)/2"
		clientTextY := "(h-th)/1.25"
		eventTextX := "(w-tw)/2"
		eventTextY := "H-th-100"
		eventTextMargin := 100

		cuttedDir := "cortado"
		editedDir := "editado"

		// If eventText is not empty, wrap it as needed
		if eventText != "" {
			eventText = wrapText(eventText, fontSize, videoWidth, eventTextMargin)
		}
		// If clientText is not empty, wrap it as needed
		if clientText != "" {
			clientText = wrapText(clientText, fontSize, videoWidth, eventTextMargin)
		}

	},
}

func init() {
	// Do not really need this i think, at least for now.
	// editCmd.Flags().String("event", "e", "", "Name of the event.")
	editCmd.Flags().StringVarP(&logo, "logo", "l", "", "Partial path to logo image file (~/Videos/eventos/assets/logo/<your_path>).")
	editCmd.Flags().StringVarP(&music, "music", "m", "", "Partial path to music file (~/Videos/eventos/assets/musica/cortado/<your_path>).")
	editCmd.Flags().StringVarP(&image, "image", "i", "", "Client image file name, file has to be inside <event-name>/.")
	editCmd.Flags().StringVarP(&eventText, "text", "t", "", "Text to overlay on the whole video.")

	editCmd.Flags().StringVarP(&clientText, "client", "c", "", "Client text to display over client image.")
	editCmd.Flags().StringVarP(&up, "up", "u", "", "Upper line of client text (manual split, bypasses auto-wrap).")
	editCmd.MarkFlagsMutuallyExclusive("client", "up")

	editCmd.Flags().StringVarP(&down, "down", "d", "", "Lower line of client text (manual split, bypasses auto-wrap).")
	editCmd.Flags().StringVarP(&clientColor, "client-color", "C", "#E6E70F", "Client text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\").")
	editCmd.Flags().StringVarP(&textColor, "text-color", "T", "#E6E70F", "Event text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\").")
	editCmd.Flags().StringVarP(&font, "font", "f", "", "Font filename to use for event text (e.g., MyFont.ttf).")
	editCmd.Flags().StringVarP(&leftIcon, "left", "L", "", "Path to icon image to display to the left of client text.")
	editCmd.Flags().StringVarP(&rightIcon, "right", "R", "", "Path to icon image to display to the right of client text.")

	editCmd.Flags().BoolVarP(&save, "save", "s", false, "Save the current command.")
	editCmd.Flags().BoolVarP(&run, "run", "r", false, "Run the saved command.")
	editCmd.MarkFlagsMutuallyExclusive("save", "run")
	// You either `run` the saved command or executed/`save` the current one, making `logo`
	// required.
	editCmd.MarkFlagsOneRequired("run", "logo")

	rootCmd.AddCommand(editCmd)
}

func wrapText(text string, fontSize int, videoW int, margin int) string {
	// Use a 0.7 factor (aprox 56px) for safer wrap.
	return wordwrap.WrapString(text, uint((videoW-margin)/(fontSize*7/10)))
}
