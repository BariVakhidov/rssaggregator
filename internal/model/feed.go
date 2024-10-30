package model

import (
	"time"

	"github.com/google/uuid"
)

type FeedInfo struct {
	Name string
	Url  string
}

type FeedFollowInfo struct {
	FeedId uuid.UUID `json:"feed_id"`
}

type FeedFollowsInfo struct {
	UserID uuid.UUID
	Limit  int
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}
