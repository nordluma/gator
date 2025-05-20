package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nordluma/gator/internal/database"
	"github.com/nordluma/gator/internal/rss"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the `login` handler expects a single argument, the username")
	}

	name := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return err
	}

	s.config.SetUser(user.Name)

	fmt.Printf("%s has been set as user\n", user.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the `register` handler expects a single argument, the username")
	}

	name := cmd.args[0]
	now := time.Now()
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	s.config.SetUser(name)

	fmt.Printf("%s was created\n", user.Name)

	return nil
}

func handlerReset(s *state, _ command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return err
	}

	fmt.Println("users have been reset")

	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerFetchFeed(_ *state, cmd command) error {
	feedUrl := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %v\n", feed)

	return nil
}
