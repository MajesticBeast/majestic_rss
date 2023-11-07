-- name: CreateFeed :exec
INSERT INTO feeds (name, url, last_post_title, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: UpdateFeedLastPostTitle :exec
UPDATE feeds SET last_post_title = $1, updated_at = $2 WHERE id = $3;