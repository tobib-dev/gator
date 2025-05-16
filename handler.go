package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	userName := sql.NullString{String: name, Valid: true}
	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		//return fmt.Errorf("Username does not exist in database, register user using '<register> name'")
		os.Exit(1)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User switched successfully to " + name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	userName := sql.NullString{String: cmd.Args[0], Valid: true}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}
	_, err := s.db.GetUser(context.Background(), userName)
	if err == nil {
		fmt.Println("User already exists!")
		os.Exit(1)
	}

	_, err = s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Error creating user")
	}
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("User: %s, was created at: %v", cmd.Args[0], userParams.CreatedAt)

	return nil
}
