package feed

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/BariVakhidov/rssaggregator/internal/storage"
	"github.com/google/uuid"
)

type FeedProvider interface {
	Feeds(ctx context.Context) ([]model.Feed, error)
	CreateFeed(ctx context.Context, feedInfo model.FeedInfo) (*model.Feed, error)
	CreateFeedFollow(ctx context.Context, userId uuid.UUID, feedId uuid.UUID) (*model.FeedFollow, error)
	DeleteFeedFollow(ctx context.Context, userId uuid.UUID, followId uuid.UUID) error
	FeedFollows(ctx context.Context, feedFollowsInfo model.FeedFollowsInfo) ([]model.FeedFollow, error)
	MarkFeedAsFetched(ctx context.Context, feedId uuid.UUID) (*model.Feed, error)
	NextFeedsToFetch(ctx context.Context, limit int) ([]model.Feed, error)
}

type Service struct {
	log          *slog.Logger
	feedProvider FeedProvider
}

func New(log *slog.Logger, feedProvider FeedProvider) *Service {
	return &Service{
		log:          log,
		feedProvider: feedProvider,
	}
}

func (s *Service) Feeds(ctx context.Context) ([]model.Feed, error) {
	const op = "service.feed.Feeds"
	log := s.log.With(slog.String("op", op))

	feeds, err := s.feedProvider.Feeds(ctx)
	if err != nil {
		log.Error("failed retrieving feeds", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feeds, nil
}

func (s *Service) CreateFeed(ctx context.Context, feedInfo model.FeedInfo) (*model.Feed, error) {
	const op = "service.feed.CreateFeed"
	log := s.log.With(slog.String("op", op))

	feed, err := s.feedProvider.CreateFeed(ctx, feedInfo)
	if err != nil {
		log.Error("create feed failed", sl.Err(err))
		if errors.Is(err, storage.ErrFeedExists) {
			return nil, fmt.Errorf("%s: %w", op, ErrFeedExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feed, nil
}

func (s *Service) CreateFeedFollow(ctx context.Context, userId, feedId uuid.UUID) (*model.FeedFollow, error) {
	const op = "service.feed.CreateFeedFollow"
	log := s.log.With(slog.String("op", op))

	feedFollow, err := s.feedProvider.CreateFeedFollow(ctx, userId, feedId)
	if err != nil {
		log.Error("create feed follow", sl.Err(err))
		if errors.Is(err, storage.ErrFeedExists) {
			return nil, fmt.Errorf("%s:%w", op, ErrFeedExists)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return feedFollow, nil
}

func (s *Service) DeleteFeedFollow(ctx context.Context, userId, followId uuid.UUID) error {
	const op = "service.feed.DeleteFeedFollow"
	log := s.log.With(slog.String("op", op))

	err := s.feedProvider.DeleteFeedFollow(ctx, userId, followId)
	if err != nil {
		log.Error("delete feed follow failed", sl.Err(err))
		if errors.Is(err, storage.ErrForeignKeyViolation) {
			return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *Service) FeedFollows(ctx context.Context, feedFollowsInfo model.FeedFollowsInfo) ([]model.FeedFollow, error) {
	const op = "service.feed.FeedFollows"
	log := u.log.With(
		slog.String("op", op),
	)

	follows, err := u.feedProvider.FeedFollows(ctx, feedFollowsInfo)
	if err != nil {
		log.Error("failed feed follows", sl.Err(err))
		if errors.Is(err, storage.ErrForeignKeyViolation) || errors.Is(err, storage.ErrInvalidParams) {
			return nil, fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return follows, nil
}

func (s *Service) MarkFeedAsFetched(ctx context.Context, feedId uuid.UUID) (*model.Feed, error) {
	const op = "service.feed.MarkFeedAsFetched"
	log := s.log.With(slog.String("op", op))

	feed, err := s.feedProvider.MarkFeedAsFetched(ctx, feedId)
	if err != nil {
		log.Error("mark feed ad fetched failed", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feed, nil
}

func (s *Service) NextFeedsToFetch(ctx context.Context, limit int) ([]model.Feed, error) {
	const op = "service.feed.NextFeedsToFetch"
	log := s.log.With(slog.String("op", op))

	feeds, err := s.feedProvider.NextFeedsToFetch(ctx, limit)
	if err != nil {
		log.Error("failed retrieving feeds", sl.Err(err))
		if errors.Is(err, storage.ErrInvalidParams) {
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feeds, nil
}
