-- name: CreateUser :one
INSERT INTO "users" (first_name, last_name, email, password)
VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email ILIKE $1;
