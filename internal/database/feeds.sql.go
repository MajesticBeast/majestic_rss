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
INSERT INTO feeds (name, feed_url, webhook_url, last_post_title, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateFeedParams struct {
	Name          string
	FeedUrl       string
	WebhookUrl    string
	LastPostTitle string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) error {
	_, err := q.db.ExecContext(ctx, createFeed,
		arg.Name,
		arg.FeedUrl,
		arg.WebhookUrl,
		arg.LastPostTitle,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const deleteFeed = `-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1
`

func (q *Queries) DeleteFeed(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteFeed, id)
	return err
}

const getFeedByID = `-- name: GetFeedByID :one
SELECT id, name, last_post_title, feed_url, webhook_url, created_at, updated_at FROM feeds WHERE id = $1
`

func (q *Queries) GetFeedByID(ctx context.Context, id int32) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByID, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.LastPostTitle,
		&i.FeedUrl,
		&i.WebhookUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFeedByName = `-- name: GetFeedByName :one
SELECT id, name, last_post_title, feed_url, webhook_url, created_at, updated_at FROM feeds WHERE name = $1
`

func (q *Queries) GetFeedByName(ctx context.Context, name string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByName, name)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.LastPostTitle,
		&i.FeedUrl,
		&i.WebhookUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT id, name, last_post_title, feed_url, webhook_url, created_at, updated_at FROM feeds
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
			&i.FeedUrl,
			&i.WebhookUrl,
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

const updateFeed = `-- name: UpdateFeed :exec
UPDATE feeds SET name = $1, feed_url = $2, webhook_url = $3, updated_at = $4 WHERE id = $5
`

type UpdateFeedParams struct {
	Name       string
	FeedUrl    string
	WebhookUrl string
	UpdatedAt  time.Time
	ID         int32
}

func (q *Queries) UpdateFeed(ctx context.Context, arg UpdateFeedParams) error {
	_, err := q.db.ExecContext(ctx, updateFeed,
		arg.Name,
		arg.FeedUrl,
		arg.WebhookUrl,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}

const updateFeedLastPostTitle = `-- name: UpdateFeedLastPostTitle :exec
UPDATE feeds SET last_post_title = $1, updated_at = $2 WHERE id = $3
`

type UpdateFeedLastPostTitleParams struct {
	LastPostTitle string
	UpdatedAt     time.Time
	ID            int32
}

func (q *Queries) UpdateFeedLastPostTitle(ctx context.Context, arg UpdateFeedLastPostTitleParams) error {
	_, err := q.db.ExecContext(ctx, updateFeedLastPostTitle, arg.LastPostTitle, arg.UpdatedAt, arg.ID)
	return err
}
