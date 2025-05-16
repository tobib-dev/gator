package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Please enter a valid username, username can only be one argument!")
	}

	userName := cmd.args[0]
	err := s.cfg.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Println("Username " + cmd.args[0] + " has been set.")

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handlerFunc, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("Command invalid, please enter a valid command!")
	}
	return handlerFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	_, ok := c.handlers[name]
	if ok {
		return fmt.Errorf("Command already registered")
	}

	c.handlers[name] = f
	return nil
}
