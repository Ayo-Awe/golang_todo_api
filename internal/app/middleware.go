package app

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

type contextKey string

const (
	defaultPerPage = 20
	maxPerPage     = 100
	defaultCursor  = 2_147_483_647
)

const userContextKey contextKey = "user"
const pagingContextKey contextKey = "paging"

func (a *Application) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			render.Render(w, r, ErrUnauthorized("Missing authorizaton header"))
			return
		}

		headerComponents := strings.Split(authHeader, " ")
		if len(headerComponents) != 2 {
			render.Render(w, r, ErrUnauthorized("Malformed authorization header"))
			return
		}

		if strings.ToLower(headerComponents[0]) != "basic" {
			render.Render(w, r, ErrUnauthorized("Invalid authentication type. Only basic authentication is allowed"))
			return
		}

		encodedCredentials := headerComponents[1]
		bytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			render.Render(w, r, ErrUnauthorized("Malformed credentials"))
			slog.Error(err.Error())
			return
		}

		credentials := strings.Split(string(bytes), ":")
		if len(credentials) != 2 {
			render.Render(w, r, ErrUnauthorized("Malformed basic auth credentials"))
			return
		}

		email := credentials[0]
		password := credentials[1]

		user, err := a.store.Users().GetUserByEmail(r.Context(), email)
		if err != nil {
			if errors.Is(err, ErrUserNotFound) {
				render.Render(w, r, ErrUnauthorized("Invalid credentials"))
				return
			}
			render.Render(w, r, ErrInternalServerError("An unexpected error occured"))
			slog.Error(err.Error())
			return
		}

		if !user.ComparePassword(password) {
			render.Render(w, r, ErrUnauthorized("Invalid credentials"))
			return
		}

		r = a.setUserCtx(r, user)

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func (a *Application) Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get paging
		rawCursor := r.URL.Query().Get("cursor")
		rawPerPage := r.URL.Query().Get("per_page")

		perPage, err := strconv.Atoi(rawPerPage)
		if err != nil || perPage <= 0 {
			perPage = defaultPerPage
		}

		if perPage > maxPerPage {
			perPage = maxPerPage
		}

		cursor, err := strconv.Atoi(rawCursor)
		if err != nil {
			cursor = defaultCursor
		}

		paging := Paging{
			Cursor:  cursor,
			PerPage: perPage,
		}

		r = a.setCtxPaging(r, paging)

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
