-- name: CreateFeed :one
INSERT INTO feeds(id,name,url)
VALUES ($1,$2,$3)
RETURNING *;

-- name: GetFeeds :many
SELECT * from feeds;

-- name: GetFeedsByUser :many
SELECT feeds.* from feeds INNER JOIN feed_follows ON feed.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY feeds.created_at;

-- name: GetNextFeedsToFetch :many
SELECT * from feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;