package converter

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreatePostInfoToDBCreatePostInfo(postInfo model.CreatePostInfo) database.CreatePostParams {
	description := pgtype.Text{Valid: false}
	if len(postInfo.Description) != 0 {
		description.String = postInfo.Description
		description.Valid = true
	}

	return database.CreatePostParams{
		ID:          postInfo.ID,
		Title:       postInfo.Title,
		Url:         postInfo.Url,
		FeedID:      postInfo.FeedID,
		PublishedAt: pgtype.Timestamp{Time: postInfo.PublishedAt, Valid: true},
		Description: description,
	}
}
