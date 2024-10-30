package feed

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/BariVakhidov/rssaggregator/internal/service/feed"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type FeedService interface {
	Feeds(ctx context.Context) ([]model.Feed, error)
	CreateFeed(ctx context.Context, feedInfo model.FeedInfo) (*model.Feed, error)
	CreateFeedFollow(ctx context.Context, userId uuid.UUID, feedId uuid.UUID) (*model.FeedFollow, error)
	DeleteFeedFollow(ctx context.Context, userId uuid.UUID, followId uuid.UUID) error
	FeedFollows(ctx context.Context, feedFollowsInfo model.FeedFollowsInfo) ([]model.FeedFollow, error)
}

type Implementation struct {
	log         *slog.Logger
	feedService FeedService
	validator   *validator.Validate
}

func New(log *slog.Logger, feedService FeedService) *Implementation {
	return &Implementation{
		log:         log,
		feedService: feedService,
		validator:   validator.New(),
	}
}

func (i *Implementation) checkFeedServiceErr(err error) error {
	if errors.Is(err, feed.ErrFeedExists) {
		return model.NewAPIErr(http.StatusConflict, ErrFeedExists)
	}

	if errors.Is(err, feed.ErrInvalidCredentials) {
		return model.NewAPIErr(http.StatusBadRequest, ErrInvalidCredentials)
	}

	if errors.Is(err, feed.ErrNotFound) {
		return model.NewAPIErr(http.StatusNotFound, ErrFeedNotFound)
	}

	return err
}
