package main

import "errors"

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func newCommands() commands {
	return commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if !ok {
		return errors.New("unknown command")
	}

	if err := handler(s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

type command struct {
	name string
	args []string
}
