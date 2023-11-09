-- +goose Up
CREATE TABLE feeds (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  last_post_title VARCHAR(225) NOT NULL,
  feed_url VARCHAR(255) UNIQUE NOT NULL,
  webhook_url VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE feeds;