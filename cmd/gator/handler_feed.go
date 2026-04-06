package main

import (
	"context"
	"fmt"
	"time"

	"gator/internal/database"

	"github.com/google/uuid"
)

func handlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("required 2 arguments, (name, url)")
	}

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.arguments[0], Url: cmd.arguments[1], UserID: user.ID})
	if err != nil {
		return fmt.Errorf("feed with this url already exists: %+v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %+v", err)
	}

	fmt.Printf("feed created: %+v\n", feed)
	fmt.Printf("following feed: %s\n", feedFollow.FeedName)
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

func handlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("required to pass in the URL with this command")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %+v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("failed to create feed follow record: %+v", err)
	}

	fmt.Printf("feed name: %s, user name: %s\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("required 1 argument: feed URL")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("failed to find feed: %+v", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %+v", err)
	}

	fmt.Printf("unfollowed %s\n", feed.Name)
	return nil
}

func handlerFollowers(s *State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("this command does not require any additional arguments")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch feed follows: %+v", err)
	}

	for _, v := range follows {
		fmt.Printf("feed: %s, user: %s\n", v.FeedName, v.UserName)
	}
	return nil
}
