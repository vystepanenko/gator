package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/vystepanenko/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("User name is required")
	}
	existingUser, err := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error while getting user: %s", err)
	}

	err = s.cfg.SetUser(existingUser.Name)
	if err != nil {
		return err
	}

	fmt.Println("User", existingUser.Name, "logged in")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("User name is required")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	u, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Println("Error creating user:", err)

		os.Exit(1)
	}

	s.cfg.SetUser(u.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println("Error resetting database:", err)
		os.Exit(1)
	}

	fmt.Println("Database reset")
	os.Exit(0)
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting users: %s", err)
	}

	for _, u := range users {
		printUser(s, u)
	}

	return nil
}

func printUser(s *state, user database.User) error {
	if user.Name == *s.cfg.CurrentUserName {
		fmt.Printf("* %s (current)\n", user.Name)
	} else {
		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}
