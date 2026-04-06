package main

import (
	"fmt"

	"gator/internal/config"
	"gator/internal/database"
)

type State struct {
	cfg *config.Config
	db  *database.Queries
}

type Command struct {
	name      string
	arguments []string
}

type Commands struct {
	handlers map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("handler does not exist")
	}

	return handler(s, cmd)
}

func (c *Commands) register(name string, handler func(*State, Command) error) {
	_, ok := c.handlers[name]
	if ok {
		panic(fmt.Errorf("handler already exists"))
	}

	c.handlers[name] = handler
}
