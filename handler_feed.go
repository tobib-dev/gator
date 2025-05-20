package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name><url>", cmd.Name)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("===================================")

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:          %s\n", feed.ID)
	fmt.Printf("* Created:     %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:     %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:        %s\n", feed.Name)
	fmt.Printf("* URL:         %s\n", feed.Url)
	fmt.Printf("* UserID:      %s\n", feed.UserID)
}
