package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bjornpagen/tiktok-video-processor/pkg/server/db"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktok/scraperapi"

	lmdb "wellquite.org/golmdb"
)

type Server struct {
	DB      *db.TikTokDB
	Scraper *scraperapi.Scraper
}

func New(dbPath, fetcherApiKey, scraperApiKey string) *Server {
	return &Server{
		DB:      db.New(dbPath),
		Scraper: scraperapi.New(scraperApiKey),
	}
}

func (s *Server) AddUsername(username string) error {
	log.Printf("Adding @%s to the database", username)
	userId, err := s.Scraper.FetchUserId(username)
	if err != nil {
		return err
	}

	// Fetch current userIds, append the new one
	userIds, err := s.DB.GetUserIDList()
	if err != nil {
		// if MDB_NOTFOUND, then create a new list
		if err == lmdb.NotFound {
			userIds = []string{userId}
			if err := s.DB.SetUserIDList(userIds); err != nil {
				return err
			}
			return nil
		}
	}

	// Check if the userId already exists
	for _, id := range userIds {
		if id == userId {
			log.Printf("User @%s already exists in the database", username)
			return nil
		}
	}

	userIds = append(userIds, userId)

	// Save the new userIds
	if err := s.DB.SetUserIDList(userIds); err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() error {
	// Open the TikTokDB
	if err := s.DB.Open(); err != nil {
		return err
	}
	defer s.DB.Close()

	// Run the server
	return s.runWithSignalHandling()
}

func (s *Server) runWithSignalHandling() error {
	// Create a channel to listen for os signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received an interrupt signal, stopping updates...")
		os.Exit(0)
	}()

	if err := s.UpdateAllDaily(); err != nil {
		log.Printf("Update failed: %s", err)
		return err
	}

	return nil
}

func (s *Server) UpdateAllDaily() error {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		if err := s.UpdateAllOnce(); err != nil {
			return err
		}

		// Wait for the next tick
		<-ticker.C
	}
}

func (s *Server) UpdateAllOnce() error {
	// Fetch all the userIds
	ids, err := s.DB.GetUserIDList()
	if err != nil {
		return err
	}

	// For all users, update them
	for _, userID := range ids {
		if err := s.Update(userID); err != nil {
			return err
		}

		user, err := s.DB.GetUser(userID)
		if err != nil {
			return err
		}

		log.Printf("Updated user @%s #%s", user.UniqueID, userID)
	}

	return nil
}
