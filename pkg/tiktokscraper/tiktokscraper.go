package tiktokscraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type TikTokScraper struct {
	APIHost   string
	APIKey    string
	RateLimit time.Duration
}

func New(apiKey string) *TikTokScraper {
	return &TikTokScraper{
		APIHost:   "tiktok-best-experience.p.rapidapi.com",
		APIKey:    apiKey,
		RateLimit: 20 * time.Millisecond,
	}
}

type Response struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	MinCursor int64   `json:"min_cursor"`
	MaxCursor int64   `json:"max_cursor"`
	HasMore   int     `json:"has_more"`
	AwemeList []Aweme `json:"aweme_list"`
}

type Aweme struct {
	AwemeID        string       `json:"aweme_id"`
	Desc           string       `json:"desc"`
	CreateTime     int64        `json:"create_time"`
	Author         Author       `json:"author"`
	Music          Music        `json:"music"`
	ChaList        []ChaList    `json:"cha_list"`
	Video          Video        `json:"video"`
	ShareURL       string       `json:"share_url"`
	Statistics     Statistics   `json:"statistics"`
	Status         Status       `json:"status"`
	Rate           int          `json:"rate"`
	TextExtra      []TextExtra  `json:"text_extra"`
	LabelTop       Image        `json:"label_top"`
	ShareInfo      ShareInfo    `json:"share_info"`
	AwemeType      int          `json:"aweme_type"`
	AuthorUserID   int64        `json:"author_user_id"`
	IsHashTag      int          `json:"is_hash_tag"`
	Region         string       `json:"region"`
	GroupID        string       `json:"group_id"`
	DescLanguage   string       `json:"desc_language"`
	MiscInfo       string       `json:"misc_info"`
	DistributeType int          `json:"distribute_type"`
	VideoControl   VideoControl `json:"video_control"`
}

type TextExtra struct {
	Start       int    `json:"start"`
	End         int    `json:"end"`
	HashtagID   string `json:"hashtag_id,omitempty"`
	HashtagName string `json:"hashtag_name,omitempty"`
	SecUID      string `json:"sec_uid,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	Type        int    `json:"type,omitempty"`
}

type VideoControl struct {
	AllowDownload         bool `json:"allow_download"`
	ShareType             int  `json:"share_type"`
	ShowProgressBar       int  `json:"show_progress_bar"`
	DraftProgressBar      int  `json:"draft_progress_bar"`
	AllowDuet             bool `json:"allow_duet"`
	AllowReact            bool `json:"allow_react"`
	PreventDownloadType   int  `json:"prevent_download_type"`
	AllowDynamicWallpaper bool `json:"allow_dynamic_wallpaper"`
	TimerStatus           int  `json:"timer_status"`
	AllowMusic            bool `json:"allow_music"`
	AllowStitch           bool `json:"allow_stitch"`
}

type Statistics struct {
	AwemeID            string `json:"aweme_id"`
	CommentCount       int    `json:"comment_count"`
	DiggCount          int    `json:"digg_count"`
	DownloadCount      int    `json:"download_count"`
	PlayCount          int    `json:"play_count"`
	ShareCount         int    `json:"share_count"`
	ForwardCount       int    `json:"forward_count"`
	LoseCount          int    `json:"lose_count"`
	LoseCommentCount   int    `json:"lose_comment_count"`
	WhatsAppShareCount int    `json:"whatsapp_share_count"`
}

type Status struct {
	AwemeID        string       `json:"aweme_id"`
	IsDelete       bool         `json:"is_delete"`
	AllowShare     bool         `json:"allow_share"`
	AllowComment   bool         `json:"allow_comment"`
	PrivateStatus  int          `json:"private_status"`
	InReviewing    bool         `json:"in_reviewing"`
	Reviewed       int          `json:"reviewed"`
	SelfSee        bool         `json:"self_see"`
	IsProhibited   bool         `json:"is_prohibited"`
	DownloadStatus int          `json:"download_status"`
	ReviewResult   ReviewResult `json:"review_result"`
	VideoMute      VideoMute    `json:"video_mute"`
}

type ReviewResult struct {
	ReviewStatus int `json:"review_status"`
}

type VideoMute struct {
	IsMute   bool   `json:"is_mute"`
	MuteDesc string `json:"mute_desc"`
}

type Author struct {
	Avatar168x168          Image     `json:"avatar_168x168"`
	Avatar300x300          Image     `json:"avatar_300x300"`
	AvatarLarger           Image     `json:"avatar_larger"`
	AvatarMedium           Image     `json:"avatar_medium"`
	AvatarThumb            Image     `json:"avatar_thumb"`
	AvatarURI              string    `json:"avatar_uri"`
	AwemeCount             int       `json:"aweme_count"`
	CoverURL               []Image   `json:"cover_url"`
	DuetSetting            int       `json:"duet_setting"`
	EnterpriseVerifyReason string    `json:"enterprise_verify_reason"`
	FavoritingCount        int       `json:"favoriting_count"`
	FollowerCount          int       `json:"follower_count"`
	FollowingCount         int       `json:"following_count"`
	InsID                  string    `json:"ins_id"`
	Language               string    `json:"language"`
	Nickname               string    `json:"nickname"`
	OriginalMusician       Musician  `json:"original_musician"`
	Region                 string    `json:"region"`
	SecUID                 string    `json:"sec_uid"`
	ShareInfo              ShareInfo `json:"share_info"`
	ShortID                string    `json:"short_id"`
	TotalFavorited         int64     `json:"total_favorited"`
	UID                    string    `json:"uid"`
	UniqueID               string    `json:"unique_id"`
	UniqueIDModifyTime     int64     `json:"unique_id_modify_time"`
	UserMode               int       `json:"user_mode"`
	UserRate               int       `json:"user_rate"`
	VerificationType       int       `json:"verification_type"`
	VideoIcon              Image     `json:"video_icon"`
}

type Music struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Album         string `json:"album"`
	IsOriginal    bool   `json:"is_original"`
	OwnerID       string `json:"owner_id"`
	OwnerNickname string `json:"owner_nickname"`
}

type ChaList struct {
	Desc    string `json:"desc"`
	Type    int    `json:"type"`
	CID     string `json:"cid"`
	ChaName string `json:"cha_name"`
}

type Video struct {
	Width        int  `json:"width"`
	Height       int  `json:"height"`
	Duration     int  `json:"duration"`
	HasWatermark bool `json:"has_watermark"`
}

type ShareInfo struct {
	ShareDescInfo string `json:"share_desc_info"`
}

type Musician struct {
	DiggCount      int `json:"digg_count"`
	MusicCount     int `json:"music_count"`
	MusicUsedCount int `json:"music_used_count"`
}

type Image struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
}

func (t *TikTokScraper) FetchUserFeedData(username string, maxCursor int64) (*Data, error) {
	url := fmt.Sprintf("https://%s/user/%s/feed", t.APIHost, username)
	if maxCursor > 0 {
		url = fmt.Sprintf("%s?max_cursor=%d", url, maxCursor)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", t.APIKey)
	req.Header.Add("X-RapidAPI-Host", t.APIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, errors.New("failed to fetch user feed")
	}

	return &response.Data, nil
}

func (t *TikTokScraper) FetchUserAwemesFromMinCursor(username string, minCursor int64) ([]Aweme, error) {
	var allAwemes []Aweme
	var maxCursor int64

	for {
		data, err := t.FetchUserFeedData(username, maxCursor)
		if err != nil {
			return nil, err
		}

		for _, aweme := range data.AwemeList {
			if aweme.CreateTime >= minCursor {
				allAwemes = append(allAwemes, aweme)
			} else {
				return allAwemes, nil
			}
		}

		if data.HasMore == 0 {
			break
		}

		maxCursor = data.MaxCursor
		time.Sleep(t.RateLimit)
	}

	return allAwemes, nil
}

func (t *TikTokScraper) FetchUserAwemes(username string) ([]Aweme, error) {
	minCursor := int64(0)
	return t.FetchUserAwemesFromMinCursor(username, minCursor)
}
