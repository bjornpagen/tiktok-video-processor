package tiktokscanner

import (
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktokdb"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktokscraper"
)

type Client struct {
	db      *tiktokdb.TikTokDB
	scraper *tiktokscraper.Client
}
