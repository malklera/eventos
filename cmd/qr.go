package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// qrCmd represents the qr command
var qrCmd = &cobra.Command{
	Use:   "qr <url>",
	Short: "Generate the QR code image for the event.",
	Long:  `e.g. ./eventos qr <url>`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(qrCmd)
}
