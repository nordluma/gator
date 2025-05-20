package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nordluma/gator/internal/config"
	"github.com/nordluma/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	s := state{db: dbQueries, config: &cfg}

	cmds := newCommands()
	cmds.register("register", handlerRegister)
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
