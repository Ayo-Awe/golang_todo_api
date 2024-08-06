package main

import (
	"fmt"
	"log"

	"github.com/ayo-awe/golang_todo_api/internal/app"
	"github.com/ayo-awe/golang_todo_api/internal/database"
)

//	@title			Task Managment API
//	@version		1.0
//	@description	This is a task management api server.

// @securityDefinitions.basic	BasicAuth
// @BasePath					/api
func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	database, err := database.New(cfg.DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApplication(cfg, database)

	if err := app.Start(); err != nil {
		fmt.Print(err)
	}
}
