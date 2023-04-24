-- name: CreateUser :one
INSERT INTO users (username, password, email, first_name)
VALUES ($1, $2, $3, $4) RETURNING *;


-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1 LIMIT 1;


-- COALESCE return first null value
-- SELECT COALESCE(NULL, NULL, 2, 'W3Schools.com'); return 2

-- using nullable argument feature of sqlc so that each fields
-- can updated independently without affecting each other


-- name: UpdateUser :one
UPDATE users
SET password            = coalesce(sqlc.narg(password), password),
    password_changed_at = coalesce(sqlc.narg(password_changed_at), password_changed_at),
    first_name          = coalesce(sqlc.narg(first_name), first_name),
    email               = coalesce(sqlc.narg(email), email),
    is_email_verified   = coalesce(sqlc.narg(is_email_verified), is_email_verified)
WHERE username = coalesce(sqlc.arg(username), username) RETURNING *;;



-- name: ListUsers :many
-- SELECT *
-- FROM authors
-- WHERE role = $1
-- ORDER BY username LIMIT $2
-- OFFSET $3;

