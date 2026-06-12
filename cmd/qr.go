/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// qrCmd represents the qr command
var qrCmd = &cobra.Command{
	Use:   "qr <url>",
	Short: "Generate the QR code image for the event.",
	Long:  `e.g. ./eventos qr <url>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 1:
			qrc, err := qrcode.New(args[0])
			if err != nil {
				return err
			}
			w, err := standard.New("qrCode.jpeg")
			if err != nil {
				return err
			}
			if err := qrc.Save(w); err != nil {
				return err
			}
		case 0:
			return fmt.Errorf("have to provide an <url>")
		default:
			return fmt.Errorf("too many arguments")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(qrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// qrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// qrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
