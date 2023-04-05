package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bjornpagen/tiktok-video-processor/pkg/server/scanner"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktok/scraperapi"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktokdb"
)

type Server struct {
	scanner *scanner.Client
	userIds []string
}

func New(dbPath, fetcherApiKey, scraperApiKey string) *Server {
	return &Server{
		scanner: scanner.New(tiktokdb.New(dbPath), scraperapi.New(scraperApiKey)),
		userIds: nil,
	}
}

func (s *Server) AddUsername(username string) error {
	userId, err := s.scanner.Scraper.FetchUserId(username)
	if err != nil {
		return err
	}

	s.userIds = append(s.userIds, userId)
	return nil
}

func (s *Server) Run() error {
	// Open the TikTokDB
	if err := s.scanner.DB.Open(); err != nil {
		return err
	}
	defer s.scanner.DB.Close()

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
	// For all users, update them
	for _, userID := range s.userIds {
		if err := s.scanner.Update(userID); err != nil {
			return err
		}

		user, err := s.scanner.DB.GetUser(userID)
		if err != nil {
			return err
		}

		log.Printf("Updated user @%s #%s", user.UniqueID, userID)
	}

	return nil
}
