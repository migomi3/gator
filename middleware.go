package main

import (
	"context"

	"github.com/migomi3/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	f := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}

		return nil
	}

	return f
}
