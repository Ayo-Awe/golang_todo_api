package database

import (
	"context"

	"github.com/ayo-awe/golang_todo_api/internal/app"
	"github.com/jackc/pgx/v5/pgxpool"
)

// database is a concrete store
type Database struct {
	conn     *pgxpool.Pool
	taskRepo app.TaskRepository
	userRepo app.UserRepository
}

func (d *Database) Users() app.UserRepository {
	return d.userRepo
}

func (d *Database) Tasks() app.TaskRepository {
	return d.taskRepo
}

func New(dsn string) (*Database, error) {
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	userRepo := NewUserRepository(conn)
	taskRepo := NewTaskRepository(conn)

	db := &Database{conn: conn, userRepo: userRepo, taskRepo: taskRepo}
	return db, nil
}
