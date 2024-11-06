package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/vystepanenko/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Name and url required")
	}

	feedParam := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	newFeed, err := s.db.CreateFeed(context.Background(), feedParam)
	if err != nil {
		return fmt.Errorf("Error creating feed: %s", err)
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(
		context.Background(),
		followParams,
	)
	if err != nil {
		return fmt.Errorf("Error while creating follow feed: %s", err)
	}

	printFeed(newFeed, user)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %s", err)
	}

	for _, f := range feeds {
		user, err := s.db.GetUser(context.Background(), f.UserID)
		if err != nil {
			return fmt.Errorf("error getting user to print feed: %s", err)
		}

		printFeed(f, user)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("URL is required")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error while getting feed: %s", err)
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(
		context.Background(),
		followParams,
	)
	if err != nil {
		return fmt.Errorf("Error while creating follow feed: %s", err)
	}

	fmt.Println("Follow feed")
	fmt.Printf("* Feed Name:    %s\n", feedFollow.FeedName)
	fmt.Printf("* User Name:    %s\n", feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)
	if err != nil {
		return err
	}

	fmt.Println("User feed's:")
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("URL is required")
	}

	feed, err := s.db.GetFeedByUrl(
		context.Background(),
		cmd.Args[0],
	)
	if err != nil {
		return fmt.Errorf("Error while getting feed: %s", err)
	}

	param := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(
		context.Background(),
		param,
	)
	if err != nil {
		return fmt.Errorf("Error while deleting feed: %s", err)
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) error {
	fmt.Println("Feed:")
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserName:      %s\n", user.Name)
	fmt.Printf("* LastFetchedAt:  %v\n", feed.LastFetchedAt.Time)
	return nil
}
