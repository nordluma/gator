package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/nordluma/gator/internal/config"
)

type state struct {
	config *config.Config
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func newCommands() commands {
	return commands{
		cmds: make(map[string]func(*state, command) error),
	}
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmds[cmd.name]
	if !ok {
		return errors.New("unknown command")
	}

	if err := handler(s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the `login` handler expects a single argument, the username")
	}

	user := cmd.args[0]
	fmt.Printf("user: %s\n", user)
	s.config.SetUser(user)
	fmt.Printf("%s has been set as user\n", user)

	return nil
}

func main() {
	cfg := config.Read()
	s := state{config: &cfg}
	cmds := newCommands()
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	if err := cmds.run(&s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
