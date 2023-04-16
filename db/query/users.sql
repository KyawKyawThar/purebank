-- name: CreateUser :one
INSERT INTO users (username, password, email, first_name)
VALUES ($1, $2, $3, $4) RETURNING *;


-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1 LIMIT 1;


-- name: ListUsers :many
-- SELECT *
-- FROM authors
-- WHERE role = $1
-- ORDER BY username LIMIT $2
-- OFFSET $3;

