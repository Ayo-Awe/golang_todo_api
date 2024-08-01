-- name: CreateTask :one
INSERT INTO "tasks" (title, description, user_id) VALUES
($1,$2,$3) RETURNING *;

-- name: GetTasks :many
SELECT * FROM "tasks"
WHERE user_id = sqlc.arg('user_id') AND id <= sqlc.arg('cursor') AND (is_completed = sqlc.narg('is_completed') OR sqlc.narg('is_completed') IS NULL)
ORDER BY id DESC
LIMIT sqlc.arg('limit');

-- name: GetTaskByID :one
SELECT * FROM "tasks"
WHERE user_id = $1 AND id = $2;

-- name: UpdateTask :one
UPDATE "tasks"
SET	title = $2,
	description = $3,
	is_completed = $4,
	updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;


-- name: DeleteTask :exec
DELETE FROM "tasks"
WHERE id = $1 AND user_id = $2;
