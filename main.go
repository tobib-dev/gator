package main

import (
	"fmt"
	"log"

	"github.com/tobib-dev/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name    string
	handler []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("tobi")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)
}
