package scanner

import (
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktok/scraperapi"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktokdb"
)

type Client struct {
	db      *tiktokdb.TikTokDB
	scraper *scraperapi.Scraper
}

func New(db *tiktokdb.TikTokDB, scraper *scraperapi.Scraper) *Client {
	return &Client{
		db:      db,
		scraper: scraper,
	}
}

func (c *Client) Update(userID string) error {
	// Refresh the user info
	user, err := c.scraper.FetchUserData(userID)
	if err != nil {
		return err
	}
	if err := c.db.SetUser(userID, user); err != nil {
		return err
	}

	// Fetch new Awemes
	awemeList := c.db.GetAwemeList(userID)
	minCursor := int64(0)
	if len(awemeList) > 0 {
		minCursor = awemeList[len(awemeList)-1].CreateTime
	}

	newAwemes, err := c.scraper.FetchUserAwemeListFromMinCursor(userID, minCursor)
	if err != nil {
		return err
	}

	// Append new Awemes to the existing list
	awemeList = append(awemeList, newAwemes...)
	if err := c.db.SetAwemeList(userID, awemeList); err != nil {
		return err
	}

	return nil
}

func (c *Client) FullUpdate(userID string) error {
	// Refresh the user info
	user, err := c.scraper.FetchUserData(userID)
	if err != nil {
		return err
	}
	if err := c.db.SetUser(userID, user); err != nil {
		return err
	}

	// Refetch all Awemes
	awemeList, err := c.scraper.FetchUserAwemeList(userID)
	if err != nil {
		return err
	}
	if err := c.db.SetAwemeList(userID, awemeList); err != nil {
		return err
	}

	return nil
}
