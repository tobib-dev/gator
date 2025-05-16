package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tobib-dev/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	gatorState := state{}
	gatorState.cfg = &cfg

	gatorCmd := commands{}
	gatorCmd.handlers = make(map[string]func(*state, command) error)
	gatorCmd.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments were provided")
		os.Exit(1)
	}
	userCmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = gatorCmd.run(&gatorState, userCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
