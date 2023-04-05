package server

import (
	lmdb "wellquite.org/golmdb"
)

func (s *Server) Update(userID string) error {
	// Refresh the user info
	user, err := s.Scraper.FetchUserData(userID)
	if err != nil {
		return err
	}
	if err := s.DB.SetUser(userID, user); err != nil {
		return err
	}

	// Fetch new Awemes
	awemeList, err := s.DB.GetAwemeList(userID)
	if err != nil {
		// if MDB_NOTFOUND, then fetch all Awemes
		if err == lmdb.NotFound {
			awemeList, err = s.Scraper.FetchUserAwemeList(userID)
			if err != nil {
				return err
			}
			if err := s.DB.SetAwemeList(userID, awemeList); err != nil {
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

	newAwemes, err := s.Scraper.FetchUserAwemeListAfterCursor(userID, minCursor)
	if err != nil {
		return err
	}

	// Append new Awemes to the existing list
	awemeList = append(awemeList, newAwemes...)
	if err := s.DB.SetAwemeList(userID, awemeList); err != nil {
		return err
	}

	return nil
}

func (s *Server) FullUpdate(userID string) error {
	// Refresh the user info
	user, err := s.Scraper.FetchUserData(userID)
	if err != nil {
		return err
	}
	if err := s.DB.SetUser(userID, user); err != nil {
		return err
	}

	// Refetch all Awemes
	awemeList, err := s.Scraper.FetchUserAwemeList(userID)
	if err != nil {
		return err
	}
	if err := s.DB.SetAwemeList(userID, awemeList); err != nil {
		return err
	}

	return nil
}
