package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/nordluma/gator/internal/config"
)

type state struct {
	config *config.Config
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
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	if err := cmds.run(&s, cmd); err != nil {
		log.Fatal(err)
	}
}
