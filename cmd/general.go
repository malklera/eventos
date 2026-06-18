package cmd

import (
	"fmt"
	"os"
	"strings"
)

// listVideos return a list of .mp4 files
func listVideos(path string) ([]os.DirEntry, error) {
	allFiles, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir(%s): %w", path, err)
	}
	var videos []os.DirEntry

	for _, file := range allFiles {
		if strings.HasSuffix(file.Name(), ".mp4") {
			videos = append(videos, file)
		}
	}
	return videos, nil
}
