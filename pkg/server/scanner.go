package server

import (
	"log"

	lmdb "wellquite.org/golmdb"
)

func (s *Server) Update(userID string) error {
	log.Printf("updating user #%s", userID)
	log.Printf("fetching user data for #%s", userID)
	// Refresh the user info
	user, err := s.Scraper.FetchUserData(userID)
	if err != nil {
		return err
	}
	if err := s.DB.SetUser(userID, user); err != nil {
		return err
	}

	log.Printf("fetching new awemes for user @%s #%s", user.UniqueID, userID)
	// Fetch new Awemes
	awemeList, err := s.DB.GetAwemeList(userID)
	if err != nil {
		if err == lmdb.NotFound {
			log.Printf("user @%s #%s not found in the database", user.UniqueID, userID)
			awemeList, err = s.Scraper.FetchUserAwemeList(userID)
			if err != nil {
				return err
			}
			if err := s.DB.SetAwemeList(userID, awemeList); err != nil {
				return err
			}

			log.Printf("updated user @%s #%s for the first time", user.UniqueID, userID)
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

	log.Printf("updated user @%s #%s", user.UniqueID, userID)

	return nil
}

func (s *Server) FullUpdate(userID string) error {
	log.Printf("performing full update for user #%s", userID)
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

	log.Printf("performed full update for user @%s #%s", user.UniqueID, userID)

	return nil
}
