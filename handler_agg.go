package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
		timeFormats := []string{"Mon Jan 2 15:04:05 MST 2006", "Mon, 02 Jan 2006 15:04:05 -0700", "2006-01-02T15:04:05Z"}

		var pubTime time.Time
		success := false
		for _, format := range timeFormats {
			pubTime, err = time.Parse(format, item.PubDate)
			if err == nil {
				success = true
				break
			}
		}

		if !success {
			log.Printf("wrong date format: %v", item.PubDate)
		} else {
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
				if pqErr, ok := err.(*pq.Error); ok {
					if pqErr.Code == "23505" && pqErr.Constraint == "posts_url_key" {
						continue
					}
				}
				fmt.Println(err)
				log.Printf("error creating post: %v", err)
			}
		}
	}
	log.Printf("Feed %s collected, %v post found", feed.Name, len(feedData.Channel.Item))
}
