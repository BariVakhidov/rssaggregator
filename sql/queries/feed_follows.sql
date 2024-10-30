-- name: CreateFeedFollow :one
INSERT INTO feed_follows(id, feed_id, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows WHERE user_id = $1 LIMIT $2;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;