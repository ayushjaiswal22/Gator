-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;
-- name: GetPostsUser :many
SELECT * from posts
ORDER BY published_at DESC 
LIMIT $1;

