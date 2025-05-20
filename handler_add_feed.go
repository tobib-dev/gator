package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get current user: %w", err)
	}
	fmt.Println(currentUser.Name)
	fmt.Println(currentUser.ID)

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't add feed to current user: %w", err)
	}

	fmt.Printf("Feed: %+v\n", feed)

	return nil
}
