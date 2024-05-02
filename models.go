package main

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIkey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		APIkey:    user.ApiKey,
	}
}
