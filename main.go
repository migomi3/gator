package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/migomi3/gator/internal/config"
	"github.com/migomi3/gator/internal/database"
)

func main() {
	c, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", c.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	s := state{
		db:  dbQueries,
		cfg: &c,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)

	input := os.Args
	if len(input) < 2 {
		log.Fatal("not enough arguments to run command")
	}

	cmd := command{
		name: input[1],
		args: input[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
