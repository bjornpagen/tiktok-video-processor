package downloader

import (
	"io"
	"net/http"
	"os"
)

// Client represents a video downloader
type Client struct {
}

// DownloadVideo downloads a video from the given URL and saves it locally
func (v *Client) DownloadVideo(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
