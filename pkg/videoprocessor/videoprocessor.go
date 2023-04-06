package videoprocessor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bjornpagen/goplay/pkg/chrome"
	"github.com/bjornpagen/tiktok-video-processor/pkg/comment"
	"github.com/bjornpagen/tiktok-video-processor/pkg/downloader"
	"github.com/bjornpagen/tiktok-video-processor/pkg/storer"
)

type VideoProcessor struct {
	VideoStorer   storer.Storer
	CommentStorer storer.Storer
	ResultStorer  storer.Storer
	tmpPath       string
}

func New(videos, comments, results storer.Storer) *VideoProcessor {
	return &VideoProcessor{
		VideoStorer:   videos,
		CommentStorer: comments,
		ResultStorer:  results,
		tmpPath:       "/tmp/videoprocessor-temporary",
	}
}

func AddTimestampToFilename(filename string) string {
	ext := filepath.Ext(filename)
	basename := filename[:len(filename)-len(ext)]
	timestamp := time.Now().UnixNano()
	newFilename := fmt.Sprintf("%s-%d%s", basename, timestamp, ext)
	return newFilename
}

func (vp *VideoProcessor) FetchVideo(mediaURL string) (string, error) {
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

	// Return the output file path
	s, err := vp.VideoStorer.Store(outFileFull)
	if err != nil {
		return "", err
	}

	return s, nil
}

func (vp *VideoProcessor) FetchComment(username, commentText string) (string, error) {
	err := func() error {
		c := &comment.CommentData{
			Username: username,
			Comment:  commentText,
		}
		cb := comment.NewCommentBuilder()

		ctx := context.Background()
		defer ctx.Done()

		chrome.Cleanup()
		err := cb.Start()
		if err != nil {
			return err
		}
		defer chrome.Cleanup()

		err = cb.UpdateComment(c)
		if err != nil {
			return err
		}
		cb.DownloadComment()
		time.Sleep(1 * time.Second)
		return nil
	}()
	if err != nil {
		return "", err
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
	destFile := AddTimestampToFilename("comment.png")
	srcPath := filepath.Join(downloadsDir, filename)
	tmpPath := filepath.Join("/tmp", destFile)

	err = os.Rename(srcPath, tmpPath)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	finPath, err := vp.CommentStorer.Store(tmpPath)
	if err != nil {
		return "", err
	}

	fmt.Println("File moved successfully.")
	return finPath, nil
}

func (vp *VideoProcessor) Combine(videoPath, commentPath string) (string, error) {
	outFile := AddTimestampToFilename("combined.mp4")

	// Create vp.path if it doesn't exist
	if _, err := os.Stat(vp.tmpPath); os.IsNotExist(err) {
		err = os.MkdirAll(vp.tmpPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	outputPath := filepath.Join(vp.tmpPath, outFile)

	// Apply the overlay using FFmpeg command-line tool and save the file locally
	scaleFactor := 5
	yPositionFactor := 0.12 // Set the desired value between 0 and 1
	xPositionFactor := 0.1  // Set the desired value between 0 and 1

	videoWidth, videoHeight, err := getVideoDimensions(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to get video dimensions: %w", err)
	}

	xPosition := int(float64(videoWidth) * xPositionFactor)
	yPosition := int(float64(videoHeight) * yPositionFactor)
	filterComplex := fmt.Sprintf("[1:v]scale=iw/%d:ih/%d[scaled];[0:v][scaled]overlay=x=%d:y=%d[out]", scaleFactor, scaleFactor, xPosition, yPosition)

	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-i", commentPath,
		"-filter_complex", filterComplex,
		"-map", "[out]",
		"-map", "0:a?",
		"-c:a", "copy",
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-y", outputPath,
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to combine video and comment: %w", err)
	}

	defer os.Remove(outputPath)

	// Return the output file path
	s, err := vp.ResultStorer.Store(outputPath)
	if err != nil {
		return "", err
	}

	return s, nil
}

func getVideoDimensions(videoPath string) (int, int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0", videoPath)

	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	dimensions := strings.Split(strings.TrimSpace(string(output)), ",")

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return 0, 0, err
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}
