package app

import (
	"context"

	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/google/uuid"
)

// TODO
type Storage interface {
	SavePendingUser(ctx context.Context, pendingUserInfo model.PendingUserInfo) (model.PendingUser, error)
	PendingUserByEmail(ctx context.Context, userEmail string) (model.PendingUser, error)
	ChangePendingUserStatus(ctx context.Context, userId uuid.UUID, status string) (model.PendingUser, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error)
	Feeds(ctx context.Context) ([]model.Feed, error)
	CreateFeed(ctx context.Context, feedInfo model.FeedInfo) (*model.Feed, error)
	CreateFeedFollow(ctx context.Context, userId uuid.UUID, feedId uuid.UUID) (*model.FeedFollow, error)
	DeleteFeedFollow(ctx context.Context, userId uuid.UUID, followId uuid.UUID) error
	FeedFollows(ctx context.Context, feedFollowsInfo model.FeedFollowsInfo) ([]model.FeedFollow, error)
	MarkFeedAsFetched(ctx context.Context, feedId uuid.UUID) (*model.Feed, error)
	NextFeedsToFetch(ctx context.Context, limit int) ([]model.Feed, error)
	Posts(ctx context.Context, userId uuid.UUID, limit int) ([]model.Post, error)
}

// type serviceProvider struct {
// 	storage Storage
// }
