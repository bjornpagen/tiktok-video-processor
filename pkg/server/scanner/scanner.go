package scanner

import (
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktok/scraperapi"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktokdb"

	lmdb "wellquite.org/golmdb"
)

type Client struct {
	DB      *tiktokdb.TikTokDB
	Scraper *scraperapi.Scraper
}

func New(db *tiktokdb.TikTokDB, scraper *scraperapi.Scraper) *Client {
	return &Client{
		DB:      db,
		Scraper: scraper,
	}
}

func (c *Client) Update(userID string) error {
	// Refresh the user info
	user, err := c.Scraper.FetchUserData(userID)
	if err != nil {
		return err
	}
	if err := c.DB.SetUser(userID, user); err != nil {
		return err
	}

	// Fetch new Awemes
	awemeList, err := c.DB.GetAwemeList(userID)
	if err != nil {
		// if MDB_NOTFOUND, then fetch all Awemes
		if err == lmdb.NotFound {
			awemeList, err = c.Scraper.FetchUserAwemeList(userID)
			if err != nil {
				return err
			}
			if err := c.DB.SetAwemeList(userID, awemeList); err != nil {
				return err
			}
			return nil
		}

		return err
	}

	minCursor := int64(0)
	if len(awemeList) > 0 {
		minCursor = awemeList[len(awemeList)-1].CreateTime
	}

	newAwemes, err := c.Scraper.FetchUserAwemeListAfterCursor(userID, minCursor)
	if err != nil {
		return err
	}

	// Append new Awemes to the existing list
	awemeList = append(awemeList, newAwemes...)
	if err := c.DB.SetAwemeList(userID, awemeList); err != nil {
		return err
	}

	return nil
}

func (c *Client) FullUpdate(userID string) error {
	// Refresh the user info
	user, err := c.Scraper.FetchUserData(userID)
	if err != nil {
		return err
	}
	if err := c.DB.SetUser(userID, user); err != nil {
		return err
	}

	// Refetch all Awemes
	awemeList, err := c.Scraper.FetchUserAwemeList(userID)
	if err != nil {
		return err
	}
	if err := c.DB.SetAwemeList(userID, awemeList); err != nil {
		return err
	}

	return nil
}
