package app

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID        int       `json:"id"`
	Firstname string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

func (u *User) SetNewPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"update_at"`
}

type TaskFilter struct {
	IsCompleted null.Bool
}

type PaginationData struct {
	NextCursor int `json:"next_cursor"`
	ItemCount  int `json:"item_count"`
	PerPage    int `json:"per_page"`
}

type Paging struct {
	Cursor  int
	PerPage int
}

type Store interface {
	Users() UserRepository
	Tasks() TaskRepository
}

type UserRepository interface {
	GetUserByEmail(context.Context, string) (*User, error)
	CreateUser(context.Context, *User) (*User, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *Task) (*Task, error)
	GetTasks(ctx context.Context, userID int, taskFilter TaskFilter, paging Paging) ([]Task, PaginationData, error)
}
