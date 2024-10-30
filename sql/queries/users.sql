-- name: CreateUser :one
INSERT INTO users(id,email,name,updated_at)
VALUES ($1,$2,$3,$4)
RETURNING *;

-- name: SavePendingUser :one
INSERT INTO pending_users(id, email, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ChangeStatusOfPendingUser :one
UPDATE pending_users
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: PendingUserByEmail :one
SELECT * FROM pending_users
WHERE email = $1;

-- name: DeletePendingUserByEmail :one
DELETE FROM pending_users
WHERE email = $1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;