package main

import (
	"database/sql"
	"fmt"
	"os"

	"gator/internal/config"
	"gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("you need to pass in the command")
		os.Exit(1)
	}

	cmd := Command{name: args[0], arguments: args[1:]}

	cnf, err := config.Read()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", cnf.Dburl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := State{cfg: &cnf, db: dbQueries}

	commands := Commands{handlers: make(map[string]func(*State, Command) error)}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerFeed)
	err = commands.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
