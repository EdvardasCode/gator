package main

import (
	"context"
	"fmt"
	"time"

	"gator/internal/database"

	"github.com/google/uuid"
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
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("required to pass in a URL")
	}

	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error: %+v", err)
	}

	fmt.Printf("%+v \n", rss)

	return nil
}
