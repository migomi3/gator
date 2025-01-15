package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("missing argument for time between reps")
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return nil
	}

	fmt.Printf("Collecting feeds every %s\n", cmd.args[0])

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("not enough arguments to run Add Feed command")
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	ffParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), ffParams)
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsWithUser(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feeds {
		fmt.Println(f)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("missing url argument")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed [%s] now followed by current user [%s]", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, f := range feeds {
		fmt.Println(f.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("missing url argument")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	return nil
}

func handlerBrowse(s *state, cmd command) error {
	var limit int32 = 2
	if len(cmd.args) >= 1 {
		parsedInt, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(parsedInt)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), limit)
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Println(p)
	}

	return nil
}
