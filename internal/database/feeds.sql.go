// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const clearFeeds = `-- name: ClearFeeds :exec
DELETE FROM feeds
`

func (q *Queries) ClearFeeds(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, clearFeeds)
	return err
}

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at, name, url, user_id
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const getFeedsWithUser = `-- name: GetFeedsWithUser :many
SELECT f.name, url, u.name AS creator
FROM feeds f
LEFT JOIN users u
ON f.user_id = u.id
`

type GetFeedsWithUserRow struct {
	Name    string
	Url     string
	Creator sql.NullString
}

func (q *Queries) GetFeedsWithUser(ctx context.Context) ([]GetFeedsWithUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedsWithUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsWithUserRow
	for rows.Next() {
		var i GetFeedsWithUserRow
		if err := rows.Scan(&i.Name, &i.Url, &i.Creator); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}