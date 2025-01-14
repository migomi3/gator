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

-- name: GetFeed :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;