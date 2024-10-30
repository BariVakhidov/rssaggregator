package model

import (
	"time"

	"github.com/google/uuid"
)

type PostInfo struct {
	UserId uuid.UUID
	Limit  int
}

type CreatePostInfo struct {
	ID          uuid.UUID
	FeedID      uuid.UUID
	Title       string
	Url         string
	PublishedAt time.Time
	Description string
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
