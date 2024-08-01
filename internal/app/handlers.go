package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"gopkg.in/guregu/null.v4"
)

func (a *Application) setUserCtx(r *http.Request, user *User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (a *Application) getCtxUser(r *http.Request) *User {
	user, ok := r.Context().Value(userContextKey).(*User)
	if !ok {
		panic("missing user in request context")
	}

	return user
}

func (a *Application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userRepo := a.store.Users()
	var payload RegiserUserRequest

	if err := render.Bind(r, &payload); err != nil {
		render.Render(w, r, ErrBadRequest("Invalid request body"))
		return
	}

	if err := payload.Validate(); err != nil {
		render.Render(w, r, ErrBadRequest(err.Error()))
		return
	}

	existingUser, err := userRepo.GetUserByEmail(r.Context(), payload.Email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	if existingUser != nil {
		render.Render(w, r, ErrConflict("Existing user email"))
		return
	}

	userPayload := &User{
		Firstname: payload.Firstname,
		Lastname:  payload.Lastname,
		Email:     payload.Email,
	}

	if err := userPayload.SetNewPassword(payload.Password); err != nil {
		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	user, err := userRepo.CreateUser(r.Context(), userPayload)
	if err != nil {
		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	res := RegisterUserResponse{*user}
	render.Render(w, r, NewSuccessResponse(res))
}

func (a *Application) CreateTask(w http.ResponseWriter, r *http.Request) {
	user := a.getCtxUser(r)

	var requestBody CreateTaskRequest
	if err := render.Bind(r, &requestBody); err != nil {
		render.Render(w, r, ErrBadRequest("Invalid request body"))
		return
	}

	if err := requestBody.Validate(); err != nil {
		render.Render(w, r, ErrBadRequest(err.Error()))
		return
	}

	taskPayload := &Task{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		UserID:      user.ID,
	}

	newTask, err := a.store.Tasks().CreateTask(r.Context(), taskPayload)
	if err != nil {
		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewSuccessResponse(CreateTaskResponse{Task: *newTask}))
}

func (a *Application) GetTasks(w http.ResponseWriter, r *http.Request) {
	user := a.getCtxUser(r)

	rawCursor := r.URL.Query().Get("cursor")
	rawPerPage := r.URL.Query().Get("per_page")
	status := r.URL.Query().Get("status")

	perPage, err := strconv.Atoi(rawPerPage)
	if err != nil {
		perPage = 20
	}

	cursor, err := strconv.Atoi(rawCursor)
	if err != nil {
		cursor = 1_000_000_000_000
	}

	var isCompleted null.Bool
	if status == "completed" {
		isCompleted = null.BoolFrom(true)
	} else if status == "pending" {
		isCompleted = null.BoolFrom(false)
	}

	paging := Paging{Cursor: cursor, PerPage: perPage}
	tasks, paginationData, err := a.store.Tasks().GetTasks(r.Context(), user.ID, TaskFilter{IsCompleted: isCompleted}, paging)
	if err != nil {
		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	payload := NewSuccessResponse(GetTasksResponse{tasks}).WithPaginationData(paginationData)
	render.Render(w, r, payload)
}
