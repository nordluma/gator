package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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

func handlerFetchFeed(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New(
			"the `agg` handler expects a single argument, the duration between intervals",
		)
	}

	interval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s\n", interval)
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	feed, err := rss.FetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:     uuid.New(),
			FeedID: feedToFetch.ID,
			Title:  item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}

			log.Printf("Couldn't create post: %v", err)
			continue
		}

	}

	log.Printf(
		"Feed %s collected, %v posts found",
		feedToFetch.Name,
		len(feed.Channel.Item),
	)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("the `addfeed` handler expects two arguments, feed name and url")
	}

	feedName := cmd.args[0]
	url := cmd.args[1]

	now := time.Now()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return nil
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf(
			"Name: %s\nUrl: %s\nUser Name: %s\n",
			feed.Name,
			feed.Url,
			feed.UserName,
		)
	}

	return nil
}

func handlerFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New(
			"the `feed` handler expects a single argument, the url of the feed to follow",
		)
	}

	feedToFollow, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	followedFeed, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:     uuid.New(),
			UserID: user.ID,
			FeedID: feedToFollow.ID,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("Followed %s\n", followedFeed.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range followedFeeds {
		fmt.Printf("%s\n", feed)
	}

	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New(
			"the `unfollow` handler expects a single argument, the url of the feed to unfollow",
		)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return nil
	}

	if err = s.db.DeleteFeedFollowForUser(
		context.Background(),
		database.DeleteFeedFollowForUserParams{
			UserID: user.ID,
			FeedID: feed.ID,
		}); err != nil {
		return nil
	}

	return nil
}
