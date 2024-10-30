package post

import (
	"context"
	"log/slog"

	"github.com/BariVakhidov/rssaggregator/internal/model"
)

type Service interface {
	Posts(ctx context.Context, postInfo model.PostInfo) ([]model.Post, error)
}

type Implementation struct {
	log         *slog.Logger
	postService Service
}

func New(log *slog.Logger, postService Service) *Implementation {
	return &Implementation{
		log:         log,
		postService: postService,
	}
}
