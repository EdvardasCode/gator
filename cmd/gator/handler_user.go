package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gator/internal/database"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerReset(s *State, cmd Command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("failed to truncate the table")
	}

	fmt.Println("database reset successfully")
	return nil
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("login command expects exactly 1 argument, the username")
	}

	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil || user.Name != cmd.arguments[0] {
		return fmt.Errorf("user does not exist")
	}

	s.cfg.SetUser(cmd.arguments[0])
	fmt.Println("user has been set")
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("login command expects exactly 1 argument, the username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), Name: cmd.arguments[0], CreatedAt: time.Now(), UpdatedAt: time.Now()})
	if err != nil {
		return err
	}

	s.cfg.SetUser(cmd.arguments[0])
	fmt.Println("user was created")
	fmt.Printf("%+v\n", user)
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	if len(cmd.arguments) > 0 {
		return fmt.Errorf("users command does not expect any extra arguments")
	}

	userNames, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch users")
	}

	for _, name := range userNames {
		if s.cfg.CurrentUserName == name {
			fmt.Printf("%+v (current)\n", name)
			continue
		}
		fmt.Printf("%+v\n", name)
	}

	return nil
}

func handlerAgg(s *State, cmd Command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("invalid duration %q: %+v", cmd.arguments[0], err)
	}

	fmt.Printf("collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

var rssTimeFormats = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC3339,
	"2006-01-02T15:04:05Z",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
}

func parsePublishedAt(raw string) sql.NullTime {
	for _, format := range rssTimeFormats {
		if t, err := time.Parse(format, raw); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	return sql.NullTime{}
}

func scrapeFeeds(s *State) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("error getting next feed: %+v\n", err)
		return
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("error marking feed as fetched: %+v\n", err)
		return
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("error fetching feed %s: %+v\n", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: parsePublishedAt(item.PubDate),
			FeedID:      feed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			fmt.Printf("error saving post %s: %+v\n", item.Link, err)
		}
	}
}
