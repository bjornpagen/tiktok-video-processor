package videoprocessor

import (
	"path/filepath"

	"github.com/bjornpagen/tiktok-video-processor/pkg/downloader"
	"github.com/bjornpagen/tiktok-video-processor/pkg/storer"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktok/fetcherapi"
)

type VideoProcessor struct {
	Fetcher *fetcherapi.Fetcher
	Storer  storer.Storer
	path    string
}

type VideoConfig struct {
	MediaURL       string
	ProfilePicture string
	Username       string
	Comment        string
	IsVerified     bool
}

func New(f *fetcherapi.Fetcher, s storer.Storer) *VideoProcessor {
	return &VideoProcessor{
		Fetcher: f,
		Storer:  s,
		path:    "/tmp/videoprocessor",
	}
}

func (vp *VideoProcessor) ProcessVideo(cfg *VideoConfig) (string, error) {
	// Download the video
	videoURL, err := vp.Fetcher.GetVideoURL(cfg.MediaURL)
	if err != nil {
		return "", err
	}

	toPath := filepath.Join(vp.path, videoURL)

	dl := downloader.New()
	dl.DownloadVideo(videoURL, toPath)

	// TODO: Generate the overlay image with cfg.ProfilePicture, cfg.Username, cfg.Comment, and cfg.IsVerified

	// Apply the overlay using an FFmpeg binding for Go and save the file locally

	// Return the output file path
	s, err := vp.Storer.Store(toPath)
	if err != nil {
		return "", err
	}

	return s, nil
}
