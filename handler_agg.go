package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBtwRequest, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBtwRequest)
	ticker := time.NewTicker(timeBtwRequest)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("couldn't retrieve next feeds", err)
		return
	}
	log.Println("Successfully fetched next feed!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("couldn't retrieve feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		pubTime, err := time.Parse("2025-05-22 11:21:35", item.PubDate)
		if err != nil {
			fmt.Printf("Error parsing time to time type: %w\n", err)
			log.Printf("couldn't parse time: %v", err)
		}

		_, err = db.CreatePosts(context.Background(), database.CreatePostsParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubTime,
			FeedID:      feed.ID,
		})
		if err != nil {
			fmt.Println(err)
			log.Printf("error creating post: %v", err)
		}
	}
	log.Printf("Feed %s collected, %v post found", feed.Name, len(feedData.Channel.Item))
}
