package post

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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type PostProvider interface {
	Posts(ctx context.Context, userId uuid.UUID, limit int) ([]model.Post, error)
}

type Service struct {
	log          *slog.Logger
	postProvider PostProvider
}

func New(log *slog.Logger, postProvider PostProvider) *Service {
	return &Service{
		log:          log,
		postProvider: postProvider,
	}
}

func (s *Service) Posts(ctx context.Context, postInfo model.PostInfo) ([]model.Post, error) {
	const op = "service.post.Posts"
	log := s.log.With(
		slog.String("op", op),
		slog.String("userId", postInfo.UserId.String()),
	)

	posts, err := s.postProvider.Posts(ctx, postInfo.UserId, postInfo.Limit)
	if err != nil {
		log.Error("failed retrieving user posts", sl.Err(err))
		if errors.Is(err, storage.ErrForeignKeyViolation) || errors.Is(err, storage.ErrInvalidParams) {
			return nil, fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return posts, nil
}
