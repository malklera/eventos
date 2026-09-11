package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// file path to the svg where to insert the qr
var template = ""

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }

// qrCmd represents the qr command
var qrCmd = &cobra.Command{
	Use:   "qr <url>",
	Short: "Generate the QR code image for the event.",
	Long:  `e.g. ./eventos qr <url>`,
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return preRunQR(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		qrc, err := qrcode.New(args[0])
		if err != nil {
			return err
		}
		if cmd.Flags().Changed("template") {
			uriImg, err := createURIImg(qrc)
			if err != nil {
				return err
			}
			svg, err := insertImg(uriImg)
			if err != nil {
				return err
			}
			if err := exportPDF(svg, "qrImprimir.pdf"); err != nil {
				return err
			}
		} else {
			w, err := standard.New("qrCode.jpeg")
			if err != nil {
				return err
			}
			if err := qrc.Save(w); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	qrCmd.Flags().StringVarP(&template, "template", "t", "", "Path to template file where to insert the qr.")
	rootCmd.AddCommand(qrCmd)
}

// preRunQR checks that the file name passed as argument exists
func preRunQR(cmd *cobra.Command) error {
	if cmd.Flags().Changed("template") {
		_, err := os.Stat(template)
		if err != nil {
			return err
		}
	}

	_, err := exec.LookPath("inkscape")
	if err != nil {
		return err
	}

	return nil
}

// createURIImg takes a QRCode struct and returns the URI image to be inserted
// in the svg
func createURIImg(qrc *qrcode.QRCode) (string, error) {
	buf := bytes.Buffer{}
	w := standard.NewWithWriter(nopWriteCloser{&buf}, standard.WithBuiltinImageEncoder(standard.PNG_FORMAT))
	defer w.Close()
	if err := qrc.Save(w); err != nil {
		return "", err
	}
	img := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + img, nil
}

// isertImg takes the uri containing the image and replace the 'replace' element
// from [template] returning the new svg
func insertImg(uri string) (*etree.Document, error) {
	svg := etree.NewDocument()
	if err := svg.ReadFromFile(template); err != nil {
		return nil, err
	}

	path := etree.MustCompilePath("//*")
	repl := svg.FindElementsPath(path)
	for i, c := range repl {
		if strings.HasPrefix(c.SelectAttrValue("id", "nil"), "replace") {
			svg.Root().CreateChild("image", func(e *etree.Element) {
				e.CreateAttr("id", fmt.Sprintf("qr%d", i))
				e.CreateAttr("x", c.SelectAttrValue("x", "nil"))
				e.CreateAttr("y", c.SelectAttrValue("y", "nil"))
				e.CreateAttr("width", c.SelectAttrValue("width", "nil"))
				e.CreateAttr("height", c.SelectAttrValue("height", "nil"))
				e.CreateAttr("preserveAspectRatio", "none")
				e.CreateAttr("xlink:href", uri)
			})
			c.Parent().RemoveChildAt(c.Index())
		}
	}
	return svg, nil
}

// exportPDF takes a SVG and saves it as a PDF using inkscape
func exportPDF(svg *etree.Document, file string) error {
	src, err := svg.WriteToString()
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp("", "qr-*.svg")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := tmp.WriteString(src); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	cmd := exec.Command("inkscape", tmpPath, "--export-type=pdf", "--export-filename="+file)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("inkscape failed: %w\n%s", err, out)
	}
	return nil
}
