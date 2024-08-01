-- name: CreateTask :one
INSERT INTO "tasks" (title, description, user_id) VALUES
($1,$2,$3) RETURNING *;

-- name: GetTasks :many
SELECT * FROM "tasks"
WHERE user_id = sqlc.arg('user_id') AND id <= sqlc.arg('cursor') AND (is_completed = sqlc.narg('is_completed') OR sqlc.narg('is_completed') IS NULL)
ORDER BY id DESC
LIMIT sqlc.arg('limit');
