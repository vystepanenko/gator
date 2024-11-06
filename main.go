package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/vystepanenko/gator/internal/config"
	"github.com/vystepanenko/gator/internal/database"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	st := state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		registeredCmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	cmd := command{
		Name: args[0],
		Args: args[1:],
	}
	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), *s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error while getting user: %s", err)
		}

		return handler(s, cmd, user)
	}
}
