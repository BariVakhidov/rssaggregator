package db

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
)

func DatabaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		APIkey:    user.ApiKey,
	}
}

func DatabaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		UpdatedAt: dbFeed.UpdatedAt,
		CreatedAt: dbFeed.CreatedAt,
		Name:      dbFeed.Name,
		UserID:    dbFeed.UserID,
		Url:       dbFeed.Url,
	}
}

func DatabaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feedsArr := make([]Feed, 0, len(dbFeeds))

	for _, dbFeed := range dbFeeds {
		feedsArr = append(feedsArr, DatabaseFeedToFeed(dbFeed))
	}

	return feedsArr
}

func DatabaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		CreatedAt: dbFeedFollow.CreatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func DatabaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollowsArr := make([]FeedFollow, 0, len(dbFeedFollows))

	for _, dbFeedFollow := range dbFeedFollows {
		feedFollowsArr = append(feedFollowsArr, DatabaseFeedFollowToFeedFollow(dbFeedFollow))
	}

	return feedFollowsArr
}

func DatabasePostToPost(dbPost database.Post) Post {
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

func DatabasePostsToPosts(dbPosts []database.Post) []Post {
	postsArr := make([]Post, 0, len(dbPosts))

	for _, dbPost := range dbPosts {
		postsArr = append(postsArr, DatabasePostToPost(dbPost))
	}

	return postsArr
}
