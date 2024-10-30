package converter

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

func DatabaseUserToUser(user database.User) model.User {
	return model.User{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt.Time,
		CreatedAt: user.CreatedAt.Time,
	}
}

func DatabasePendingUserToPendingUser(dbPendingUser database.PendingUser) model.PendingUser {
	return model.PendingUser{
		ID:        dbPendingUser.ID,
		UpdatedAt: dbPendingUser.UpdatedAt.Time,
		CreatedAt: dbPendingUser.CreatedAt.Time,
		Email:     dbPendingUser.Email,
		Name:      dbPendingUser.Name,
		Status:    dbPendingUser.Status,
	}
}

func DatabasePendingUsersToPendingUsers(dbPendingUsers []database.PendingUser) []model.PendingUser {
	users := make([]model.PendingUser, len(dbPendingUsers))
	for i, dbUser := range dbPendingUsers {
		users[i] = DatabasePendingUserToPendingUser(dbUser)
	}

	return users
}
