-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id) 
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    inserted.*,
    u.name AS user_name,
    f.name AS feed_name
FROM 
    inserted
JOIN 
    users u ON inserted.user_id = u.id
JOIN 
    feeds f ON inserted.feed_id = f.id;


-- name: GetFeedFollowsForUser :many
select
	ff.*,
	u."name" as user_name,
	f."name" as feed_name
from
	feed_follows ff
join users u on
	u.id = ff.user_id
join feeds f on
	f.id = ff.feed_id
where
	ff.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows ff
where ff.feed_id = $1 and user_id = $2;
