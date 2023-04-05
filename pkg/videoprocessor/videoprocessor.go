package videoprocessor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bjornpagen/tiktok-video-processor/pkg/downloader"
	"github.com/bjornpagen/tiktok-video-processor/pkg/storer"
)

type VideoProcessor struct {
	Storer storer.Storer
	path   string
}

func New(s storer.Storer) *VideoProcessor {
	return &VideoProcessor{
		Storer: s,
		path:   "/tmp/videoprocessor-temporary",
	}
}

func (vp *VideoProcessor) ProcessVideo(mediaURL string) (string, error) {
	outFile := AddTimestampToFilename("video.mp4")

	// Create vp.path if it doesn't exist
	if _, err := os.Stat(vp.path); os.IsNotExist(err) {
		err = os.MkdirAll(vp.path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	outFileFull := filepath.Join(vp.path, outFile)

	dl := downloader.New()
	dl.DownloadVideo(mediaURL, outFileFull)
	defer os.Remove(outFileFull)

	// TODO: Generate the overlay image with cfg.ProfilePicture, cfg.Username, cfg.Comment, and cfg.IsVerified

	// Apply the overlay using an FFmpeg binding for Go and save the file locally

	// Return the output file path
	s, err := vp.Storer.Store(outFileFull)
	if err != nil {
		return "", err
	}

	return s, nil
}

func AddTimestampToFilename(filename string) string {
	ext := filepath.Ext(filename)
	basename := filename[:len(filename)-len(ext)]
	timestamp := time.Now().UnixNano()
	newFilename := fmt.Sprintf("%s-%d%s", basename, timestamp, ext)
	return newFilename
}
