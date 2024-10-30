-- name: CreatePost :one
INSERT INTO posts( id, title, url, feed_id, published_at, description )
VALUES ( $1, $2, $3, $4, $5, $6 )
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE url = $1;


-- name: GetPostsForUser :many
SELECT posts.* from posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;