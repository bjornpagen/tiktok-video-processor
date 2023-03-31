package tiktokfetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/ratelimit"
)

type TikTokFetcher struct {
	APIHost    string
	APIKey     string
	RateLimit  ratelimit.Limiter
	HttpClient *http.Client
}

func New(apiKey string) *TikTokFetcher {
	return &TikTokFetcher{
		APIHost:   "tiktok-download-without-watermark.p.rapidapi.com",
		APIKey:    apiKey,
		RateLimit: ratelimit.New(60),
		HttpClient: &http.Client{
			Timeout: 10,
		},
	}
}

type Response struct {
	Code          int     `json:"code"`
	Msg           string  `json:"msg"`
	ProcessedTime float64 `json:"processed_time"`
	Data          *Data   `json:"data,omitempty"`
}

type Data struct {
	AwemeID       string    `json:"aweme_id"`
	ID            string    `json:"id"`
	Region        string    `json:"region"`
	Title         string    `json:"title"`
	Cover         string    `json:"cover"`
	OriginCover   string    `json:"origin_cover"`
	Duration      int       `json:"duration"`
	Play          string    `json:"play"`
	WmPlay        string    `json:"wmplay"`
	Size          int       `json:"size"`
	WmSize        int       `json:"wm_size"`
	Music         string    `json:"music"`
	MusicInfo     MusicInfo `json:"music_info"`
	PlayCount     int       `json:"play_count"`
	DiggCount     int       `json:"digg_count"`
	CommentCount  int       `json:"comment_count"`
	ShareCount    int       `json:"share_count"`
	DownloadCount int       `json:"download_count"`
	CreateTime    int64     `json:"create_time"`
	Author        Author    `json:"author"`
}

type MusicInfo struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Play     string `json:"play"`
	Cover    string `json:"cover"`
	Author   string `json:"author"`
	Original bool   `json:"original"`
	Duration int    `json:"duration"`
	Album    string `json:"album"`
}

type Author struct {
	ID       string `json:"id"`
	UniqueID string `json:"unique_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// GetVideoURL fetches the video URL using the unofficial TikTok API
func (t *TikTokFetcher) GetVideoURL(tiktokURL string) (string, error) {
	encodedURL := url.QueryEscape(tiktokURL)
	apiURL := fmt.Sprintf("https://%s/analysis?url=%s&hd=1", t.APIHost, encodedURL)

	req, _ := http.NewRequest("GET", apiURL, nil)

	req.Header.Add("X-RapidAPI-Key", t.APIKey)
	req.Header.Add("X-RapidAPI-Host", t.APIHost)

	t.RateLimit.Take()
	res, err := t.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response Response
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
