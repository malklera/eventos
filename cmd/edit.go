/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		videoW := 1080
		// videoH := 1920
		transitionT := 1
		clientImgT := 2

		// Text styling
		// lineSpacing := 10
		fontSize := 80
		// fontColor "#E6E70F"
		// borderW := 2
		// borderColor := "#000000"
		// boxPadding := 8
		// clientTextX := "(w-tw)/2"
		// clientTextY := "(h-th)/1.25"
		// eventTextX := "(w-tw)/2"
		// eventTextY := "H-th-100"
		eventTextMargin := 100

		cuttedDir := "cortado"
		// editedDir := "editado"

		// If eventText is not empty, wrap it as needed
		if eventText != "" {
			eventText = wrapText(eventText, fontSize, videoW, eventTextMargin)
		}
		// If clientText is not empty, wrap it as needed
		if clientText != "" {
			clientText = wrapText(clientText, fontSize, videoW, eventTextMargin)
		}

		// TODO: check that all arg passed exists

		cuttedVideos, err := listVideos(cuttedDir)
		if err != nil {
			return fmt.Errorf("listVideos(%s): %w", cuttedDir, err)
		}

		editedCount := 0
		fmt.Println("Total videos:", len(cuttedVideos))
		for _, file := range cuttedVideos {
			videoPath := filepath.Join(cuttedDir, file.Name())
			videoDur, err := videoDuration(videoPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting the duration of '%s': %v", videoPath, err)
				continue
			}
			logoDur, err := videoDuration(logo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting the duration of '%s': %v", logo, err)
				continue
			}

			// client image + main video + logo
			if image != "" {
				preLogoVisualDur := videoDur + clientImgT
				totalDur := preLogoVisualDur + logoDur
				clientFadeOutStart := clientImgT - transitionT
				videoSlideOutStart := videoDur - transitionT
				videoEndFadeStart := preLogoVisualDur - transitionT
				segmentY := (videoDur - transitionT) / 8
				pLen := (videoDur + 2*transitionT) / 5
				part1Start := 0
				part1End := pLen

				// xfade 1 setup: part1 and part2 overlap by transitionDuration
				part2Start := part1End - transitionT
				part2End := part2Start + pLen

				part3Start := part2End
				part3End := part3Start + pLen

				part4Start := part3End
				part4End := part4Start + pLen

				part5Start := part4End - transitionT
				part5End := videoDur

				textFadeInStart := clientImgT - transitionT
				textFadeInEnd := clientImgT
				textFadeOutStart := videoEndFadeStart
				fmt.Printf("Client: %d, Video: %d, Logo: %d, Total: %d\n", clientImgT, videoDur, logoDur)

			} else {

			}
		}
		fmt.Println("Edited videos:", editedCount)
		return nil
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
	editCmd.Flags().StringVarP(&font, "font", "f", "/usr/local/share/fonts/Courgette-Regular.ttf", "Path to font to use for all text (e.g., /usr/share/fonts/TTF/MyFont.ttf).")
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

// videoDuration take a path and run ffprobe to learn its duration, return the quotient and any error
func videoDuration(path string) (int, error) {
	ffprobe := exec.Command("ffprobe",
		"-loglevel",
		"error",
		"-show_entries",
		"format=duration",
		"-output_format",
		"default=noprint_wrappers=1:nokey=1",
		path)
	ffprobe.Stderr = os.Stderr
	out, err := ffprobe.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe.Run(%s): %w", ffprobe.String(), err)
	}
	str, _, _ := strings.Cut(string(out), ".")
	dur, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi(string(%v)): %w", out, err)
	}
	return dur, nil
}

// checkArgs ensure all passed arguments exists
func checkArgs(logoPath string) error {
	_, err := os.Stat(logoPath)
	if err != nil {
		return err
	}
}
