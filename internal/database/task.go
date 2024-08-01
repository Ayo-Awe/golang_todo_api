package database

import (
	"context"

	"github.com/ayo-awe/golang_todo_api/internal/app"
	"github.com/ayo-awe/golang_todo_api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepo struct {
	queries *sqlc.Queries
	conn    *pgxpool.Pool
}

func NewTaskRepository(conn *pgxpool.Pool) app.TaskRepository {
	return &taskRepo{conn: conn, queries: sqlc.New(conn)}
}

func (repo *taskRepo) CreateTask(ctx context.Context, task *app.Task) (*app.Task, error) {
	arg := sqlc.CreateTaskParams{
		Title:       task.Title,
		Description: task.Description,
		UserID:      int32(task.UserID),
	}

	sqlcTask, err := repo.queries.CreateTask(ctx, arg)
	if err != nil {
		return nil, err
	}

	return repo.toAppTask(&sqlcTask), nil
}

func (repo *taskRepo) GetTasks(ctx context.Context, userID int, filter app.TaskFilter, paging app.Paging) ([]app.Task, app.PaginationData, error) {
	return []app.Task{}, app.PaginationData{}, nil
}

func (repo *taskRepo) toAppTask(sqlcTask *sqlc.Task) *app.Task {
	return &app.Task{
		ID:          int(sqlcTask.ID),
		Title:       sqlcTask.Title,
		Description: sqlcTask.Description,
		IsCompleted: sqlcTask.IsCompleted,
		UserID:      int(sqlcTask.UserID),
		CreatedAt:   sqlcTask.CreatedAt.Time,
		UpdatedAt:   sqlcTask.UpdatedAt.Time,
	}
}
