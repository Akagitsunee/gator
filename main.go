package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/akagitsunee/gator/internal/config"
	"github.com/akagitsunee/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	cfg    *config.Config
	dbConn *sql.DB
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Reading config failed: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Connection to database failed: %v", err)
	}
	dbQueries := database.New(db)

	s := &state{
		db:     dbQueries,
		cfg:    &cfg,
		dbConn: db,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerCreateFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments provided")
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
