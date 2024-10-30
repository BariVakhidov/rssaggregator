package converter

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

func ServicePendingUserInfoToDBParams(pendingUserInfo model.PendingUserInfo) database.SavePendingUserParams {
	return database.SavePendingUserParams{
		ID:    pendingUserInfo.ID,
		Email: pendingUserInfo.Email,
		Name:  pendingUserInfo.Name,
	}
}
