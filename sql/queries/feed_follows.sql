-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, user_id, feed_id)
    VALUES ($1, $2, $3)
    RETURNING *
)

SELECT
    inserted.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted
INNER JOIN users ON inserted.user_id = users.id
INNER JOIN feeds ON inserted.feed_id = feeds.id;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.*,
    u.name AS user_name,
    f.name AS feed_name
FROM feed_follows AS ff
INNER JOIN users AS u ON ff.user_id = u.id
INNER JOIN feeds AS f ON ff.feed_id = f.id
WHERE u.id = $1;
