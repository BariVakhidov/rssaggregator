package converter

import (
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

func DatabasePostToPost(dbPost database.Post) model.Post {
	var description *string

	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return model.Post{
		ID:          dbPost.ID,
		UpdatedAt:   dbPost.UpdatedAt.Time,
		CreatedAt:   dbPost.CreatedAt.Time,
		FeedID:      dbPost.FeedID,
		Description: description,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		PublishedAt: dbPost.PublishedAt.Time,
	}
}

func DatabasePostsToPosts(dbPosts []database.Post) []model.Post {
	postsArr := make([]model.Post, 0, len(dbPosts))

	for _, dbPost := range dbPosts {
		postsArr = append(postsArr, DatabasePostToPost(dbPost))
	}

	return postsArr
}
