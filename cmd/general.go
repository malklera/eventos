package cmd

import (
	"fmt"
	"io"
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

// copyFile copy files from src to dst
func copyFile(src string, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("os.Open(%s): %v", src, err)
	}
	defer srcF.Close()

	dstF, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("os.Create(%s): %v", dst, err)
	}
	defer dstF.Close()

	_, err = io.Copy(dstF, srcF)
	if err != nil {
		return fmt.Errorf("io.Copy(%s, %s): %v", srcF.Name(), dstF.Name(), err)
	}
	if err := dstF.Sync(); err != nil {
		return fmt.Errorf("dstF.Sync(): %v", err)
	}
	return nil
}
