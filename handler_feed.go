package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("===================================")

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("===================================")
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:          %s\n", feed.ID)
	fmt.Printf("* Created:     %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:     %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:        %s\n", feed.Name)
	fmt.Printf("* URL:         %s\n", feed.Url)
	fmt.Printf("* User:        %s\n", user.Name)
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var postLimit int32
	if len(cmd.Args) == 1 {
		limInt, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("couldn't convert string to int: %w", err)
		}
		postLimit = int32(limInt)
	} else {
		postLimit = 2
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  postLimit,
	})
	if err != nil {
		fmt.Println("couldn't get posts")
		fmt.Println(err)
		log.Printf("couldn't fetch users posts: %v", err)
	}

	for _, post := range posts {
		fmt.Println(post.Title)
	}
	return nil
}
