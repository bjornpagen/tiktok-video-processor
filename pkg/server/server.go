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
	dbPath  string
	userIds []string
	scraper *scraperapi.Scraper
}

func New(dbPath, fetcherApiKey, scraperApiKey string) *Server {
	return &Server{
		dbPath:  dbPath,
		userIds: nil,
		scraper: scraperapi.New(scraperApiKey),
	}
}

func (s *Server) AddUsername(username string) error {
	userId, err := s.scraper.FetchUserId(username)
	if err != nil {
		return err
	}

	s.userIds = append(s.userIds, userId)
	return nil
}

func (s *Server) Run() error {
	// Create a new TikTokDB
	db := tiktokdb.New(s.dbPath)

	// Open the TikTokDB
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	// Create a new Scanner
	scanner := scanner.New(db, s.scraper)

	// Run the server
	return s.runWithSignalHandling(scanner)
}

func (s *Server) runWithSignalHandling(scanner *scanner.Client) error {
	// Create a channel to listen for os signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received an interrupt signal, stopping updates...")
		os.Exit(0)
	}()

	if err := s.UpdateAllDaily(scanner); err != nil {
		log.Printf("Update failed: %s", err)
		return err
	}

	return nil
}

func (s *Server) UpdateAllDaily(scanner *scanner.Client) error {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		if err := s.UpdateAllOnce(scanner); err != nil {
			return err
		}

		// Wait for the next tick
		<-ticker.C
	}
}

func (s *Server) UpdateAllOnce(scanner *scanner.Client) error {
	// For all users, update them
	for _, userID := range s.userIds {
		if err := scanner.Update(userID); err != nil {
			return err
		}

		log.Printf("Updated user %s", userID)
	}

	return nil
}
