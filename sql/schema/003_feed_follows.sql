-- +goose Up
CREATE TABLE feed_follows(
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_id, feed_id)
);

-- +goose Down
TRUNCATE TABLE feed_follows;
