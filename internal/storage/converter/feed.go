package converter

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

func FeedFollowsInfoToDBFeedFollowsInfo(feedFollowsInfo model.FeedFollowsInfo) database.GetFeedFollowsParams {
	return database.GetFeedFollowsParams{
		UserID: feedFollowsInfo.UserID,
		Limit:  int32(feedFollowsInfo.Limit),
	}
}
