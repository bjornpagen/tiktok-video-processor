package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/bjornpagen/tiktok-video-processor/pkg/metadata"
	"github.com/bjornpagen/tiktok-video-processor/pkg/server/db"
	"github.com/bjornpagen/tiktok-video-processor/pkg/storer"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktok/fetcherapi"
	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktok/scraperapi"
	"github.com/bjornpagen/tiktok-video-processor/pkg/videoprocessor"

	lmdb "wellquite.org/golmdb"
)

type Server struct {
	DB             *db.TikTokDB
	Scraper        *scraperapi.Scraper
	Fetcher        *fetcherapi.Fetcher
	VideoStorage   storer.Storer
	CommentStorage storer.Storer
	ResultStorage  storer.Storer
}

func New(dbPath, outPath, fetcherApiKey, scraperApiKey string) *Server {
	return &Server{
		DB:             db.New(dbPath),
		Scraper:        scraperapi.New(scraperApiKey),
		Fetcher:        fetcherapi.New(fetcherApiKey),
		VideoStorage:   storer.NewLocalStorer(filepath.Join(outPath, "videos")),
		CommentStorage: storer.NewLocalStorer(filepath.Join(outPath, "comments")),
		ResultStorage:  storer.NewLocalStorer(filepath.Join(outPath, "results")),
	}
}

func (s *Server) AddUsername(username string) error {
	log.Printf("adding @%s to the database", username)
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
			log.Printf("user @%s already exists in the database", username)
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

func (s *Server) RemoveUsername(username string) error {
	userId, err := s.Scraper.FetchUserId(username)
	if err != nil {
		return err
	}

	// Fetch current userIds, remove the one to be deleted
	userIds, err := s.DB.GetUserIDList()
	if err != nil {
		return err
	}

	// Check if the userId already exists
	for i, id := range userIds {
		if id == userId {
			userIds = append(userIds[:i], userIds[i+1:]...)
			break
		}
	}

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
		log.Println("received an interrupt signal, stopping updates...")
		os.Exit(0)
	}()

	if err := s.UpdateAllDaily(); err != nil {
		log.Printf("update failed: %s", err)
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
		if err := s.FullUpdate(userID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GenerateCommentedVideo(a *scraperapi.Aweme, commentUsername, commentText, imagePath string) (string, error) {
	dlUrl, err := s.Fetcher.GetVideoURL(a.ShareURL)
	if err != nil {
		return "", err
	}

	vp := videoprocessor.New(s.VideoStorage, s.CommentStorage, s.ResultStorage)
	videoPath, err := vp.FetchVideo(dlUrl)
	if err != nil {
		return "", err
	}

	// TODO: crop video

	// TODO: change contrast and colors

	// TODO: strip audio

	// Fetch the comment
	commentPath, err := vp.FetchComment(commentUsername, commentText, imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to fetch comment: %w", err)
	}

	// Combine the video and comment
	finalPath, err := vp.Combine(videoPath, commentPath)
	if err != nil {
		return "", fmt.Errorf("failed to combine video and comment: %w", err)
	}

	// Edit metadata
	metadata.GenerateMetadataAndWriteToFile(finalPath)

	return finalPath, nil
}

func (s *Server) FetchVideo(a *scraperapi.Aweme) (string, error) {
	dlUrl, err := s.Fetcher.GetVideoURL(a.ShareURL)
	if err != nil {
		return "", err
	}

	vp := videoprocessor.New(s.VideoStorage, s.CommentStorage, s.ResultStorage)
	videoPath, err := vp.FetchVideo(dlUrl)
	if err != nil {
		return "", err
	}

	// Edit metadata
	metadata.GenerateMetadataAndWriteToFile(videoPath)

	return videoPath, nil
}

func (s *Server) FetchAllVideos(userID string) error {
	// Fetch all the awemes for the user
	awemes, err := s.DB.GetAwemeList(userID)
	if err != nil {
		return err
	}

	// Let's do it concurrently instead.
	var wg sync.WaitGroup
	for _, a := range awemes {
		wg.Add(1)
		go func(a scraperapi.Aweme) {
			defer wg.Done()
			_, err := s.FetchVideo(&a)
			if err != nil {
				log.Printf("failed to fetch video: %s", err)
			}
		}(a)
	}
	wg.Wait()

	return nil
}
