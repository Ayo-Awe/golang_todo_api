package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	config *Config
	store  Store
}

func NewApplication(config *Config, store Store) *Application {
	return &Application{config: config, store: store}
}

func (a *Application) buildRoutes() http.Handler {
	r := chi.NewRouter()
	api := chi.NewRouter()

	api.Route("/auth", func(r chi.Router) {
		r.Post("/signup", a.RegisterUser)
	})

	api.Route("/tasks", func(r chi.Router) {
		r.Use(a.basicAuthMiddleware)
		r.Post("/", a.CreateTask)
		r.With(a.Paginate).Get("/", a.GetTasks)
	})

	r.Mount("/api", api)

	return r
}

func (a *Application) Start() error {
	srv := http.Server{
		Handler: a.buildRoutes(),
		Addr:    fmt.Sprintf(":%d", a.config.PORT),
	}

	go func() {
		fmt.Printf("Starting server on port %d\n", a.config.PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	<-gracefulShutdown
	fmt.Println("Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	fmt.Println("Graceful shutdown successful...")
	return nil
}
