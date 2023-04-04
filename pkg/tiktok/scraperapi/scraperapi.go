package scraperapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/ratelimit"
)

type Scraper struct {
	APIHost    string
	APIKey     string
	RateLimit  ratelimit.Limiter
	HttpClient *http.Client
}

func New(apiKey string) *Scraper {
	return &Scraper{
		APIHost:   "tiktok-best-experience.p.rapidapi.com",
		APIKey:    apiKey,
		RateLimit: ratelimit.New(50),
		HttpClient: &http.Client{
			Timeout: 10,
		},
	}
}

type UserDataResponse struct {
	Status string     `json:"status"`
	Data   UserLookup `json:"data"`
}

type UserLookup struct {
	StatusCode int  `json:"status_code"`
	User       User `json:"user"`
}

type User struct {
	AccountType            int              `json:"account_type"`
	Avatar168x168          Avatar           `json:"avatar_168x168"`
	Avatar300x300          Avatar           `json:"avatar_300x300"`
	AvatarLarger           Avatar           `json:"avatar_larger"`
	AvatarMedium           Avatar           `json:"avatar_medium"`
	AvatarThumb            Avatar           `json:"avatar_thumb"`
	AwemeCount             int              `json:"aweme_count"`
	BioSecureURL           string           `json:"bio_secure_url"`
	BioURL                 string           `json:"bio_url"`
	CommerceUserInfo       CommerceUserInfo `json:"commerce_user_info"`
	EnterpriseVerifyReason string           `json:"enterprise_verify_reason"`
	FollowerCount          int              `json:"follower_count"`
	FollowingCount         int              `json:"following_count"`
	InsID                  string           `json:"ins_id"`
	Nickname               string           `json:"nickname"`
	OriginalMusician       OriginalMusician `json:"original_musician"`
	PrivacySetting         PrivacySetting   `json:"privacy_setting"`
	SecUID                 string           `json:"sec_uid"`
	ShareInfo              DataShareInfo    `json:"share_info"`
	ShortID                string           `json:"short_id"`
	SignatureLanguage      string           `json:"signature_language"`
	TabSettings            TabSettings      `json:"tab_settings"`
	TotalFavorited         int              `json:"total_favorited"`
	UID                    string           `json:"uid"`
	UniqueID               string           `json:"unique_id"`
	VideoIcon              VideoIcon        `json:"video_icon"`
}

type Avatar struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
}

type CommerceUserInfo struct {
	AdExperienceEntry bool     `json:"ad_experience_entry"`
	AdExperienceText  string   `json:"ad_experience_text"`
	AdRevenueRits     []string `json:"ad_revenue_rits"`
}

type OriginalMusician struct {
	DiggCount      int `json:"digg_count"`
	MusicCount     int `json:"music_count"`
	MusicUsedCount int `json:"music_used_count"`
}

type PrivacySetting struct {
	FollowingVisibility int `json:"following_visibility"`
}

type DataShareInfo struct {
	ShareURL         string `json:"share_url"`
	ShareDesc        string `json:"share_desc"`
	ShareTitle       string `json:"share_title"`
	BoolPersist      int    `json:"bool_persist"`
	ShareTitleMyself string `json:"share_title_myself"`
	ShareTitleOther  string `json:"share_title_other"`
	ShareDescInfo    string `json:"share_desc_info"`
}

type TabSettings struct {
	PrivateTab PrivateTab `json:"private_tab"`
}

type PrivateTab struct {
	PrivateTabStyle int  `json:"private_tab_style"`
	ShowPrivateTab  bool `json:"show_private_tab"`
}

type VideoIcon struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
}

type UserFeedResponse struct {
	Status string    `json:"status"`
	Data   FeedChunk `json:"data"`
}

type FeedChunk struct {
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

func (t *Scraper) FetchUserData(userId string) (*User, error) {
	url := fmt.Sprintf("https://%s/user/id/%s", t.APIHost, userId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", t.APIKey)
	req.Header.Add("X-RapidAPI-Host", t.APIHost)

	t.RateLimit.Take()
	res, err := t.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response UserDataResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %w", err)
	}

	if response.Status != "ok" {
		return nil, errors.New("endpoint failed to return ok status")
	}

	if response.Data.StatusCode != 0 {
		return nil, errors.New("endpoint failed to find user")
	}

	return &response.Data.User, nil
}

func (t *Scraper) FetchUserId(username string) (string, error) {
	url := fmt.Sprintf("https://%s/user/%s", t.APIHost, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-RapidAPI-Key", t.APIKey)
	req.Header.Add("X-RapidAPI-Host", t.APIHost)

	t.RateLimit.Take()
	res, err := t.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response UserDataResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling JSON response: %w", err)
	}

	if response.Status != "ok" {
		return "", errors.New("endpoint failed to return ok status")
	}

	if response.Data.StatusCode != 0 {
		return "", errors.New("endpoint failed to find user")
	}

	return response.Data.User.UID, nil
}

func (t *Scraper) FetchUserFeed(userId string, maxCursor int64) (*FeedChunk, error) {
	url := fmt.Sprintf("https://%s/user/id/%s/feed", t.APIHost, userId)
	if maxCursor > 0 {
		url = fmt.Sprintf("%s?max_cursor=%d", url, maxCursor)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", t.APIKey)
	req.Header.Add("X-RapidAPI-Host", t.APIHost)

	t.RateLimit.Take()
	res, err := t.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response UserFeedResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %w", err)
	}

	if response.Status != "ok" {
		return nil, errors.New("failed to fetch user feed")
	}

	return &response.Data, nil
}

func (t *Scraper) FetchUserAwemeListFromMinCursor(userId string, minCursor int64) ([]Aweme, error) {
	var allAwemes []Aweme
	var maxCursor int64

	for {
		data, err := t.FetchUserFeed(userId, maxCursor)
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
	}

	return allAwemes, nil
}

func (t *Scraper) FetchUserAwemeList(userId string) ([]Aweme, error) {
	minCursor := int64(0)
	return t.FetchUserAwemeListFromMinCursor(userId, minCursor)
}
