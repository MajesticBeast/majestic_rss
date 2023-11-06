-- +goose Up
CREATE TABLE feeds (
  id INTEGER PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  last_post_title VARCHAR(225) NOT NULL,
  url VARCHAR(255) UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE feeds;