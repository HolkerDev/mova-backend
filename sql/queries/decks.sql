-- name: GetDeckByID :one
SELECT * FROM decks WHERE id = $1;

-- name: GetDecksByUserID :many
SELECT * FROM decks WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateDeck :one
INSERT INTO decks (user_id, name, source_language, target_language)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteDeck :exec
DELETE FROM decks WHERE id = $1;
