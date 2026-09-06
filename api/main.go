package main

import (
	"log"

	"github.com/vonrimez/TaskAPI/internal/common/config"
	"github.com/vonrimez/TaskAPI/internal/database"
	"github.com/vonrimez/TaskAPI/internal/handlers"
	"github.com/vonrimez/TaskAPI/internal/service"
	"github.com/vonrimez/TaskAPI/server"
)

func main() {

	cfg := config.LoadConfig()

	log.Println("Starting server")
	log.Println("Connecting with database...")

	db, err := database.EstablishConnection("pgx", cfg.DBURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	log.Println("Server successfully connected with database")

	taskRepo := database.NewTaskDB(db)
	userRepo := database.NewUserDB(db)

	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)

	validate := service.NewValidate()

	hdl := handlers.NewHandler(
		taskService, userService, validate,
	)

	server := server.NewServer(cfg.PORT)
	server.InitRouting(hdl)

	log.Println("Running server")

	err = server.Run()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Server was shitting down")
}
