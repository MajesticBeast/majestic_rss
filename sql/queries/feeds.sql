-- name: CreateFeed :exec
INSERT INTO feeds (name, url, last_post_title, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: UpdateFeedLastPostTitle :exec
UPDATE feeds SET last_post_title = ?, updated_at = ? WHERE id = ?;