package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserStatusPending   = "pending"
	UserStatusFailed    = "failed"
	UserStatusProcessed = "processed"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserEvent struct {
	ID    uuid.UUID
	Email string
}

type PendingUser struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Name      string
	Status    string
}

type PendingUserInfo struct {
	ID    uuid.UUID
	Email string
	Name  string
}
