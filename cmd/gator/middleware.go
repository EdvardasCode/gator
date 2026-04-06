package main

import (
	"context"
	"fmt"

	"gator/internal/database"
)

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("current user is not logged in: %+v", err)
		}

		return handler(s, cmd, user)
	}
}
