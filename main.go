package main

import (
	"log"
	"os"

	"github.com/nordluma/gator/internal/config"
)

type state struct {
	config *config.Config
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
