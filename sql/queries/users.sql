-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByClerkID :one
SELECT * FROM users WHERE clerk_id = $1;

-- name: CreateUser :one
INSERT INTO users (clerk_id, email) VALUES ($1, $2) RETURNING *;
