package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handlerFunc, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("Command invalid, please enter a valid command!")
	}
	return handlerFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
