package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	userName := cmd.Args[0]
	err := s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User switched successfully to " + userName)

	return nil
}
