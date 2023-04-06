package videoprocessor

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/bjornpagen/goplay/pkg/chrome"
	"github.com/bjornpagen/tiktok-video-processor/pkg/comment"
	"github.com/bjornpagen/tiktok-video-processor/pkg/downloader"
	"github.com/bjornpagen/tiktok-video-processor/pkg/storer"
)

type VideoProcessor struct {
	Storer  storer.Storer
	tmpPath string
}

func New(s storer.Storer) *VideoProcessor {
	return &VideoProcessor{
		Storer:  s,
		tmpPath: "/tmp/videoprocessor-temporary",
	}
}

func (vp *VideoProcessor) ProcessVideo(mediaURL string) (string, error) {
	outFile := AddTimestampToFilename("video.mp4")

	// Create vp.path if it doesn't exist
	if _, err := os.Stat(vp.tmpPath); os.IsNotExist(err) {
		err = os.MkdirAll(vp.tmpPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	outFileFull := filepath.Join(vp.tmpPath, outFile)

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

func GenerateOverlayComment(username, commentText string, storer storer.Storer) (string, error) {
	destFile := AddTimestampToFilename("comment.png")
	c := &comment.CommentData{
		Username: username,
		Comment:  commentText,
	}
	cb := comment.NewCommentBuilder()
	ctx := context.Background()
	defer ctx.Done()
	err := cb.Start()
	if err != nil {
		return "", err
	}

	err = cb.UpdateComment(c)
	if err != nil {
		return "", err
	}
	cb.DownloadComment()
	time.Sleep(1 * time.Second)
	for i := 0; i < 100; i++ {
		chrome.Cleanup()
	}

	// get the current user
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	// locate the Chrome downloads directory
	downloadsDir := filepath.Join(u.HomeDir, "Downloads")
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		return "", err
	}

	// specify the file to move
	filename := "Comment.png"
	srcPath := filepath.Join(downloadsDir, filename)
	tmpPath := filepath.Join("/tmp", destFile)

	err = os.Rename(srcPath, tmpPath)
	if err != nil {
		return "", err
	}

	finPath, err := storer.Store(tmpPath)
	if err != nil {
		return "", err
	}

	fmt.Println("File moved successfully.")
	return finPath, nil
}
