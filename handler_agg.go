package main

import (
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: <time_between_reqs>")
	}
	timeBtwRequest, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't get the time between request: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBtwRequest)
	ticker := time.NewTicker(timeBtwRequest)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
