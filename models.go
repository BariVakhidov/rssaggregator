package main

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIkey    string    `json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		APIkey:    user.ApiKey,
	}
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		UpdatedAt: dbFeed.UpdatedAt,
		CreatedAt: dbFeed.CreatedAt,
		Name:      dbFeed.Name,
		UserID:    dbFeed.UserID,
		Url:       dbFeed.Url,
	}
}

func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	var feedsArr []Feed

	for _, dbFeed := range dbFeeds {
		feedsArr = append(feedsArr, databaseFeedToFeed(dbFeed))
	}

	return feedsArr
}
