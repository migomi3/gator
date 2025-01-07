package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/migomi3/gator/internal/config"
	"github.com/migomi3/gator/internal/database"
)

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expecting username")
	}

	u, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	config.SetUser(s.cfg, u.Name)

	fmt.Println("User set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expecting name")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	u, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("User Created: ")
	fmt.Printf(" ID: %s\n Created At: %v\n Updated At: %v\n Name: %s\n", u.ID, u.CreatedAt, u.UpdatedAt, u.Name)

	handlerLogin(s, cmd)

	return nil
}

func handlerReset(s *state, cmd command) error {
	return s.db.ClearUsers(context.Background())
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
			continue
		}

		fmt.Printf("* %s\n", u.Name)
	}

	return nil
}
