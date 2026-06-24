package cmd

import (
	"testing"
)

func TestVideoDuration(t *testing.T) {
	var tests = []struct {
		path    string
		outWant int
		errWant error
	}{
		{"/home/malklera/Videos/eventos/bunker/cortado/1.mp4", 29, nil},
		{"/home/malklera/Videos/eventos/bunker/cortado/2.mp4", 34, nil},
		{"/home/malklera/Videos/eventos/bunker/cortado/3.mp4", 37, nil},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			out, err := fileDuration(tt.path)
			if out != tt.outWant || err != tt.errWant {
				t.Errorf("got '%d' and '%v'\nwant '%d' and '%v'", out, err, tt.outWant, tt.errWant)
			}
		})
	}
}
