package database

import (
	"context"
	"errors"

	"github.com/ayo-awe/golang_todo_api/internal/app"
	"github.com/ayo-awe/golang_todo_api/internal/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	queries *sqlc.Queries
	conn    *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) app.UserRepository {
	queries := sqlc.New(conn)
	return &userRepo{queries: queries, conn: conn}
}

func (repo *userRepo) CreateUser(ctx context.Context, user *app.User) (*app.User, error) {
	sqlcUser, err := repo.queries.CreateUser(ctx, sqlc.CreateUserParams{
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
	})

	if err != nil {
		return nil, err
	}

	return repo.toAppUser(sqlcUser), nil
}

func (repo *userRepo) toAppUser(sqlcUser sqlc.User) *app.User {
	return &app.User{
		ID:        int(sqlcUser.ID),
		Firstname: sqlcUser.FirstName,
		Lastname:  sqlcUser.LastName,
		Email:     sqlcUser.Email,
		Password:  sqlcUser.Password,
		CreatedAt: sqlcUser.CreatedAt.Time,
		UpdatedAt: sqlcUser.UpdatedAt.Time,
	}
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string) (*app.User, error) {
	sqlcUser, err := repo.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app.ErrUserNotFound
		}
		return nil, err
	}

	return repo.toAppUser(sqlcUser), nil
}
