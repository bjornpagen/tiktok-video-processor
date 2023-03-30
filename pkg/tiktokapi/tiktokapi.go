package tiktokapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// TikTokAPI represents the TikTok API configuration
type TikTokAPI struct {
	APIHost string
	APIKey  string
}

// NewTikTokAPI returns a new TikTokAPI instance
func NewTikTokAPI(apiKey string) *TikTokAPI {
	return &TikTokAPI{
		APIHost: "tiktok-download-without-watermark.p.rapidapi.com",
		APIKey:  apiKey,
	}
}

// GetVideoURL fetches the video URL using the unofficial TikTok API
func (t *TikTokAPI) GetVideoURL(videoID string) (string, error) {
	url := fmt.Sprintf("https://%s/analysis?url=https://vm.tiktok.com/%s/&hd=0", t.APIHost, videoID)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", t.APIKey)
	req.Header.Add("X-RapidAPI-Host", t.APIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// TODO: Parse the video URL from the response body
	videoURL := ""

	return videoURL, nil
}
