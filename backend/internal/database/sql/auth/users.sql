-- name: CreateUser :one
-- model: User
-- @typedef: auth.users User
INSERT INTO auth.users (first_name, last_name, email, password)
VALUES (@first_name, @last_name, lower(@email), @password)
RETURNING *;

-- name: GetUserByEmail :one
-- @typedef: auth.users User
SELECT * FROM auth.users
WHERE email = lower(@email);

-- name: ListUsers :many
-- @typedef: auth.users User
SELECT * FROM auth.users
ORDER BY created_at DESC;
