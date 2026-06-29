package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestPreRun(t *testing.T) {
	var tests = []struct {
		name    string
		flags   map[string]string
		wantErr error
	}{
		{name: "client", flags: map[string]string{"client": "foo"}, wantErr: ErrClientWithoutImage},
		{name: "up", flags: map[string]string{"up": "foo"}, wantErr: ErrUpWithoutImage},
		{name: "left", flags: map[string]string{"left": "foo"}, wantErr: ErrLeftWithoutClient},
		{name: "right", flags: map[string]string{"right": "foo"}, wantErr: ErrRightWithoutClient},
		{name: "client-color", flags: map[string]string{"client-color": "foo"}, wantErr: ErrClientColorWithoutClient},
		{name: "text-color", flags: map[string]string{"text-color": "foo"}, wantErr: ErrTextColorWithoutText},
		{name: "font", flags: map[string]string{"font": "foo"}, wantErr: ErrFontWithoutTextOrClient},
		// {name: "", flags: map[string]string{"": "foo"}, wantErr: Err},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String("image", "", "")
			cmd.Flags().String("client", "", "")
			cmd.Flags().String("up", "", "")
			cmd.Flags().String("left", "", "")
			cmd.Flags().String("right", "", "")
			cmd.Flags().String("client-color", "", "")
			cmd.Flags().String("text", "", "")
			cmd.Flags().String("text-color", "", "")
			cmd.Flags().String("font", "", "")
			for k, v := range tt.flags {
				cmd.Flags().Set(k, v)
			}

			err := preRun(cmd)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got: '%v', want: '%v'", err, tt.wantErr)
			}
		})
	}
}

// This only works if the videos files it point to exist, change to float too
// func TestVideoDuration(t *testing.T) {
// 	var tests = []struct {
// 		path    string
// 		outWant int
// 		errWant error
// 	}{
// 		{"/home/malklera/Videos/eventos/bunker/cortado/1.mp4", 29, nil},
// 		{"/home/malklera/Videos/eventos/bunker/cortado/2.mp4", 34, nil},
// 		{"/home/malklera/Videos/eventos/bunker/cortado/3.mp4", 37, nil},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.path, func(t *testing.T) {
// 			out, err := fileDuration(tt.path)
// 			if out != tt.outWant || err != tt.errWant {
// 				t.Errorf("got '%d' and '%v'\nwant '%d' and '%v'", out, err, tt.outWant, tt.errWant)
// 			}
// 		})
// 	}
// }
