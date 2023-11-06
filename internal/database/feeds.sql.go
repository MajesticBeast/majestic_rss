// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: feeds.sql

package database

import (
	"context"
	"time"
)

const createFeed = `-- name: CreateFeed :exec
INSERT INTO feeds (name, url, last_post_title, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
`

type CreateFeedParams struct {
	Name          string
	Url           string
	LastPostTitle string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) error {
	_, err := q.db.ExecContext(ctx, createFeed,
		arg.Name,
		arg.Url,
		arg.LastPostTitle,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getFeeds = `-- name: GetFeeds :many
SELECT id, name, last_post_title, url, created_at, updated_at FROM feeds
`

func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.LastPostTitle,
			&i.Url,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
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

const updateFeedLastPostTitle = `-- name: UpdateFeedLastPostTitle :exec
UPDATE feeds SET last_post_title = ?, updated_at = ? WHERE id = ?
`

type UpdateFeedLastPostTitleParams struct {
	LastPostTitle string
	UpdatedAt     time.Time
	ID            int64
}

func (q *Queries) UpdateFeedLastPostTitle(ctx context.Context, arg UpdateFeedLastPostTitleParams) error {
	_, err := q.db.ExecContext(ctx, updateFeedLastPostTitle, arg.LastPostTitle, arg.UpdatedAt, arg.ID)
	return err
}
