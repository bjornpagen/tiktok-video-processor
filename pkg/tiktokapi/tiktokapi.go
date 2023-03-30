package tiktokapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type TikTokResponse struct {
	Code          int                 `json:"code"`
	Msg           string              `json:"msg"`
	ProcessedTime float64             `json:"processed_time"`
	Data          *TikTokResponseData `json:"data,omitempty"`
}

type TikTokResponseData struct {
	AwemeID       string          `json:"aweme_id"`
	ID            string          `json:"id"`
	Region        string          `json:"region"`
	Title         string          `json:"title"`
	Cover         string          `json:"cover"`
	OriginCover   string          `json:"origin_cover"`
	Duration      int             `json:"duration"`
	Play          string          `json:"play"`
	WmPlay        string          `json:"wmplay"`
	Size          int             `json:"size"`
	WmSize        int             `json:"wm_size"`
	Music         string          `json:"music"`
	MusicInfo     TikTokMusicInfo `json:"music_info"`
	PlayCount     int             `json:"play_count"`
	DiggCount     int             `json:"digg_count"`
	CommentCount  int             `json:"comment_count"`
	ShareCount    int             `json:"share_count"`
	DownloadCount int             `json:"download_count"`
	CreateTime    int64           `json:"create_time"`
	Author        TikTokAuthor    `json:"author"`
}

type TikTokMusicInfo struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Play     string `json:"play"`
	Cover    string `json:"cover"`
	Author   string `json:"author"`
	Original bool   `json:"original"`
	Duration int    `json:"duration"`
	Album    string `json:"album"`
}

type TikTokAuthor struct {
	ID       string `json:"id"`
	UniqueID string `json:"unique_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// GetVideoURL fetches the video URL using the unofficial TikTok API
func (t *TikTokAPI) GetVideoURL(tiktokURL string) (string, error) {
	encodedURL := url.QueryEscape(tiktokURL)
	apiURL := fmt.Sprintf("https://%s/analysis?url=%s&hd=1", t.APIHost, encodedURL)

	req, _ := http.NewRequest("GET", apiURL, nil)

	req.Header.Add("X-RapidAPI-Key", t.APIKey)
	req.Header.Add("X-RapidAPI-Host", t.APIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response TikTokResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Code != 0 {
		return "", fmt.Errorf("API error: %s", response.Msg)
	}

	videoURL := response.Data.Play
	return videoURL, nil
}
