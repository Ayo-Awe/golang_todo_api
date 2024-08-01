-- name: CreateTask :one
INSERT INTO "tasks" (title, description, user_id) VALUES
($1,$2,$3) RETURNING *;
