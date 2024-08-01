package app

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ErrorResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

type SuccessResponse struct {
	Status string          `json:"status"`
	Data   interface{}     `json:"data"`
	Paging *PaginationData `json:"paging,omitempty"`
}

func (s *SuccessResponse) WithPaginationData(p PaginationData) *SuccessResponse {
	s.Paging = &p
	return s
}

func (s *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Status: "success",
		Data:   data,
	}
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func ErrBadRequest(msg string) render.Renderer {
	return &ErrorResponse{
		Status:     "error",
		Message:    msg,
		StatusCode: http.StatusBadRequest,
	}
}

func ErrUnauthorized(msg string) render.Renderer {
	return &ErrorResponse{
		Status:     "error",
		Message:    msg,
		StatusCode: http.StatusUnauthorized,
	}
}

func ErrForbidden(msg string) render.Renderer {
	return &ErrorResponse{
		Status:     "error",
		Message:    msg,
		StatusCode: http.StatusForbidden,
	}
}

func ErrConflict(msg string) render.Renderer {
	return &ErrorResponse{
		Status:     "error",
		Message:    msg,
		StatusCode: http.StatusConflict,
	}
}

func ErrInternalServerError(msg string) render.Renderer {
	return &ErrorResponse{
		Status:     "error",
		Message:    msg,
		StatusCode: http.StatusInternalServerError,
	}
}

type RegiserUserRequest struct {
	Email     string `json:"email"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Password  string `json:"password"`
}

func (ru *RegiserUserRequest) Bind(r *http.Request) error {
	return nil
}

func (ru *RegiserUserRequest) Validate() error {
	caser := cases.Title(language.English)

	ru.Email = strings.TrimSpace(strings.ToLower(ru.Email))
	ru.Firstname = strings.TrimSpace(caser.String(ru.Firstname))
	ru.Lastname = strings.TrimSpace(caser.String(ru.Lastname))

	return validation.ValidateStruct(ru,
		validation.Field(&ru.Email, validation.Required, is.EmailFormat),
		validation.Field(&ru.Firstname, validation.Required, is.Alphanumeric),
		validation.Field(&ru.Lastname, validation.Required, is.Alphanumeric),
		validation.Field(&ru.Password, validation.Required, validation.Length(8, 255)),
	)
}

type RegisterUserResponse struct {
	User User `json:"user"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (c *CreateTaskRequest) Bind(r *http.Request) error { return nil }

func (c *CreateTaskRequest) Validate() error {
	c.Title = strings.TrimSpace(c.Title)
	c.Description = strings.TrimSpace(c.Description)

	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required),
	)
}

type CreateTaskResponse struct {
	Task Task `json:"task"`
}

type GetTasksResponse struct {
	Tasks []Task `json:"tasks"`
}
