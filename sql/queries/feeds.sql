-- name: CreateFeed :exec
INSERT INTO feeds (name, feed_url, webhook_url, last_post_title, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: UpdateFeedLastPostTitle :exec
UPDATE feeds SET last_post_title = $1, updated_at = $2 WHERE id = $3;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1;

-- name: UpdateFeed :exec
UPDATE feeds SET name = $1, feed_url = $2, webhook_url = $3, updated_at = $4 WHERE id = $5;

-- name: GetFeedByID :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetFeedByName :one
SELECT * FROM feeds WHERE name = $1;