-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows(
        id, user_id, feed_id
    ) VALUES ($1, $2, $3) 
    RETURNING *
)
SELECT 
    inserted.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted
JOIN users ON inserted.user_id = users.id
JOIN feeds ON inserted.feed_id = feeds.id;
