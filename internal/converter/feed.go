package converter

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

func DatabaseFeedToFeed(dbFeed database.Feed) model.Feed {
	return model.Feed{
		ID:        dbFeed.ID,
		UpdatedAt: dbFeed.UpdatedAt.Time,
		CreatedAt: dbFeed.CreatedAt.Time,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
	}
}

func DatabaseFeedsToFeeds(dbFeeds []database.Feed) []model.Feed {
	feedsArr := make([]model.Feed, 0, len(dbFeeds))

	for _, dbFeed := range dbFeeds {
		feedsArr = append(feedsArr, DatabaseFeedToFeed(dbFeed))
	}

	return feedsArr
}

func DatabaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) model.FeedFollow {
	return model.FeedFollow{
		ID:        dbFeedFollow.ID,
		UpdatedAt: dbFeedFollow.UpdatedAt.Time,
		CreatedAt: dbFeedFollow.CreatedAt.Time,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func DatabaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []model.FeedFollow {
	feedFollowsArr := make([]model.FeedFollow, 0, len(dbFeedFollows))

	for _, dbFeedFollow := range dbFeedFollows {
		feedFollowsArr = append(feedFollowsArr, DatabaseFeedFollowToFeedFollow(dbFeedFollow))
	}

	return feedFollowsArr
}
