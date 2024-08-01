package database

import (
	"context"
	"errors"

	"github.com/ayo-awe/golang_todo_api/internal/app"
	"github.com/ayo-awe/golang_todo_api/internal/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/guregu/null.v4"
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
	arg := sqlc.GetTasksParams{
		UserID:      int32(userID),
		Cursor:      int32(paging.Cursor),
		Limit:       int32(paging.Limit()),
		IsCompleted: pgtype.Bool(filter.IsCompleted.NullBool),
	}

	sqlcTasks, err := repo.queries.GetTasks(ctx, arg)
	if err != nil {
		return nil, app.PaginationData{}, err
	}

	tasks := make([]app.Task, len(sqlcTasks))
	for i, sqlcTask := range sqlcTasks {
		tasks[i] = *repo.toAppTask(&sqlcTask)
	}

	var nextCursor null.Int
	if len(tasks) == paging.Limit() {
		lastIndex := len(tasks) - 1
		last := tasks[lastIndex]

		nextCursor = null.IntFrom(int64(last.ID))
		tasks = tasks[:lastIndex]
	}

	paginationData := app.PaginationData{
		NextCursor: nextCursor,
		ItemCount:  len(tasks),
		PerPage:    paging.PerPage,
	}

	return tasks, paginationData, nil
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

func (repo *taskRepo) GetTaskByID(ctx context.Context, userID int, taskID int) (*app.Task, error) {
	arg := sqlc.GetTaskByIDParams{
		UserID: int32(userID),
		ID:     int32(taskID),
	}

	sqlcTask, err := repo.queries.GetTaskByID(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app.ErrTaskNotFound
		}
		return nil, err
	}

	return repo.toAppTask(&sqlcTask), nil
}

func (repo *taskRepo) UpdateTask(ctx context.Context, task *app.Task) (*app.Task, error) {
	arg := sqlc.UpdateTaskParams{
		ID:          int32(task.ID),
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
	}

	sqlcTask, err := repo.queries.UpdateTask(ctx, arg)
	if err != nil {
		return nil, err
	}

	return repo.toAppTask(&sqlcTask), nil

}
