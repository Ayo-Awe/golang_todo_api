package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (a *Application) setCtxPaging(r *http.Request, paging Paging) *http.Request {
	ctx := context.WithValue(r.Context(), pagingContextKey, paging)
	return r.WithContext(ctx)
}

func (a *Application) getCtxPaging(r *http.Request) Paging {
	paging, ok := r.Context().Value(pagingContextKey).(Paging)
	if !ok {
		panic("missing paging in request context")
	}

	return paging
}

// @Summary	Sign up
// @Tags		Auth
// @Param		request	body		RegisterUserRequest	true	"request body"
// @Success	200		{object}	SuccessResponse{data=RegisterUserResponse}
// @Failure	400,409	{object}	ErrorResponse
// @Router		/auth/signup [post]
func (a *Application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userRepo := a.store.Users()
	var payload RegisterUserRequest

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

// @Summary	Create Task
// @Tags		Tasks
// @Id			CreateTasks
// @Param		request	body		CreateTaskRequest	true	"request body"
// @Success	201		{object}	SuccessResponse{data=CreateTaskResponse}
// @Failure	400,401	{object}	ErrorResponse
// @Security	ApiKeyAuth
// @Router		/tasks [post]
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

// @Summary	Get Tasks
// @Tags		Tasks
// @Id			GetTasks
// @Param		cursor		query		int		false	"cursor for forward pagination"
// @Param		per_page	query		int		false	"maximum number of tasks to return"
// @Param		status		query		string	false	"filter by task status"	Enums(completed, pending)
// @Success	201			{object}	SuccessResponse{data=CreateTaskResponse,paging=PaginationData}
// @Failure	400,401		{object}	ErrorResponse
// @Security	BasicAuth
// @Router		/tasks [get]
func (a *Application) GetTasks(w http.ResponseWriter, r *http.Request) {
	user := a.getCtxUser(r)
	paging := a.getCtxPaging(r)

	status := r.URL.Query().Get("status")
	var isCompleted null.Bool
	if status == "completed" {
		isCompleted = null.BoolFrom(true)
	} else if status == "pending" {
		isCompleted = null.BoolFrom(false)
	}

	tasks, paginationData, err := a.store.Tasks().GetTasks(r.Context(), user.ID, TaskFilter{IsCompleted: isCompleted}, paging)
	if err != nil {
		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	payload := NewSuccessResponse(GetTasksResponse{tasks}).WithPaginationData(paginationData)
	render.Render(w, r, payload)
}

// @Summary	Edit Tasks
// @Tags		Tasks
// @Id			EditTasks
// @Param		id			path		int					true	"task id"
// @Param		request		body		EditTaskResponse	true	"request body"
// @Success	200			{object}	SuccessResponse{data=EditTaskResponse}
// @Failure	400,401,404	{object}	ErrorResponse
// @Security	BasicAuth
// @Router		/tasks/{id} [patch]
func (a *Application) EditTask(w http.ResponseWriter, r *http.Request) {
	user := a.getCtxUser(r)
	rawID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		render.Render(w, r, ErrResourceNotFound("Task not found"))
		return
	}

	task, err := a.store.Tasks().GetTaskByID(r.Context(), user.ID, id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			render.Render(w, r, ErrResourceNotFound("Task not found"))
			return
		}

		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	var requestBody EditTaskRequest
	if err := render.Bind(r, &requestBody); err != nil {
		render.Render(w, r, ErrBadRequest("Invalid request body"))
		return
	}

	if err := requestBody.Validate(); err != nil {
		render.Render(w, r, ErrBadRequest(err.Error()))
		return
	}

	if requestBody.Title != nil {
		task.Title = *requestBody.Title
	}

	if requestBody.Description != nil {
		task.Description = *requestBody.Description
	}

	if requestBody.IsCompleted != nil {
		task.IsCompleted = *requestBody.IsCompleted
	}

	updatedTask, err := a.store.Tasks().UpdateTask(r.Context(), task)
	if err != nil {
		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	render.Render(w, r, NewSuccessResponse(EditTaskResponse{*updatedTask}))
}

// @Summary	Delete Tasks
// @Tags		Tasks
// @Id			DeleteTasks
// @Param		id	path	int	true	"task id"
// @Success	204
// @Failure	401,404	{object}	ErrorResponse
// @Security	BasicAuth
// @Router		/tasks/{id} [delete]
func (a *Application) DeleteTask(w http.ResponseWriter, r *http.Request) {
	user := a.getCtxUser(r)
	rawID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		render.Render(w, r, ErrResourceNotFound("Task not found"))
		return
	}

	_, err = a.store.Tasks().GetTaskByID(r.Context(), user.ID, id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			render.Render(w, r, ErrResourceNotFound("Task not found"))
			return
		}

		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	if err := a.store.Tasks().DeleteTask(r.Context(), user.ID, id); err != nil {
		render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
		slog.Error(err.Error())
		return
	}

	render.NoContent(w, r)
}
