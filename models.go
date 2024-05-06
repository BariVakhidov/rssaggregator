package main

import (
	"time"

	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/google/uuid"
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

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FeedID      uuid.UUID `json:"feed_id"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	Description *string   `json:"description"`
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
	feedsArr := make([]Feed, 0, len(dbFeeds))

	for _, dbFeed := range dbFeeds {
		feedsArr = append(feedsArr, databaseFeedToFeed(dbFeed))
	}

	return feedsArr
}

func databaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		CreatedAt: dbFeedFollow.CreatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollowsArr := make([]FeedFollow, 0, len(dbFeedFollows))

	for _, dbFeedFollow := range dbFeedFollows {
		feedFollowsArr = append(feedFollowsArr, databaseFeedFollowToFeedFollow(dbFeedFollow))
	}

	return feedFollowsArr
}

func databasePostToPost(dbPost database.Post) Post {
	var description *string

	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		UpdatedAt:   dbPost.UpdatedAt,
		CreatedAt:   dbPost.CreatedAt,
		FeedID:      dbPost.FeedID,
		Description: description,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		PublishedAt: dbPost.PublishedAt,
	}
}

func databasePostsToPosts(dbPosts []database.Post) []Post {
	postsArr := make([]Post, 0, len(dbPosts))

	for _, dbPost := range dbPosts {
		postsArr = append(postsArr, databasePostToPost(dbPost))
	}

	return postsArr
}
