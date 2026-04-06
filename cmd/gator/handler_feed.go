package main

import (
	"context"
	"fmt"
	"time"

	"gator/internal/database"

	"github.com/google/uuid"
)

func handlerAddFeed(s *State, cmd Command) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("required 2 arguments, (name, url)")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user not found: %+v", err)
	}

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.arguments[0], Url: cmd.arguments[1], UserID: user.ID})
	if err != nil {
		return fmt.Errorf("feed with this url already exists: %+v", err)
	}

	fmt.Printf("feed created: %+v\n", feed)
	return nil
}

func handlerFeed(s *State, cmd Command) error {
	if len(cmd.arguments) > 0 {
		return fmt.Errorf("no arguments expected")
	}

	rows, err := s.db.FetchFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %+v", err)
	}

	for _, feed := range rows {
		fmt.Printf("%+v\n", feed)
	}

	return nil
}
