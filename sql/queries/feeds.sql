-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ClearFeeds :exec
DELETE FROM feeds;

-- name: GetFeedsWithUser :many
SELECT f.name, url, u.name AS creator
FROM feeds f
LEFT JOIN users u
ON f.user_id = u.id;